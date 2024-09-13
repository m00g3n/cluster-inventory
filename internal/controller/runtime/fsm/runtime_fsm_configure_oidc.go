package fsm

import (
	"context"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureOidc(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure OIDC state")

	createOpenIDConnectResource(s.instance)

	//done 0. plug in to the state machine without actually doing anything
	//TODO: - check condition if going into this state is necessary
	//         - extension enabled
	//         - OIDC config present
	//TODO: - revisit the default OIDC requirements, implementation drafted in https://github.com/kyma-project/infrastructure-manager/commit/32f016d38dafdee0ff1f9e772183825c703f046d
	//done - double check RBACs if they're in the right context (right client shoot vs gardener vs kcp )
	//done - create OpenID resource on the shoot
	//TODO: - openID resource using non-fake values
	//TODO: - (???) remove existing "additionalOidcs" from the Shoot spec
	//TODO: - recreate "additionalOidcs" in the Shoot spec
	//TODO: - extract rest client to a separate interface
	//TODO:   consider if we should stop in case of a failure
	//TODO:  the condition set should reflect framefrog condition strategy
	//TODO: - unit test the function

	//s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonGardenerCRCreated, "Unknown", "Gardener Cluster CR created, waiting for readiness")
	//return updateStatusAndRequeueAfter(controlPlaneRequeueDuration)

	// wait section - is it needed?

	//prepare subresource client to request admin kubeconfig
	srscClient := m.ShootClient.SubResource("adminkubeconfig")
	shootAdminClient, err := GetShootClient(ctx, srscClient, s.shoot)

	if err != nil {
		//TODO: handle error before updating status
		return updateStatusAndStopWithError(err)
	}

	openIDConnectResource := createOpenIDConnectResource(s.instance)
	errResourceCreation := shootAdminClient.Create(ctx, openIDConnectResource)

	if errResourceCreation != nil {
		m.log.Error(errResourceCreation, "Failed to create OpenIDConnect resource")
	}

	s.instance.UpdateStateReady(
		imv1.ConditionTypeOidcConfigured,
		imv1.ConditionReasonConfigurationCompleted,
		"OIDC configuration completed",
	)

	m.log.Info("OIDC has been configured", "Name", s.shoot.Name)

	return updateStatusAndStop()
}

func createOpenIDConnectResource(runtime imv1.Runtime) *authenticationv1alpha1.OpenIDConnect {
	//additionalOidcConfig := (*runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig)[0]

	//cr := &authenticationv1alpha1.OpenIDConnect{
	//	Spec: authenticationv1alpha1.OIDCAuthenticationSpec{
	//		IssuerURL:      *additionalOidcConfig.IssuerURL,
	//		ClientID:       *additionalOidcConfig.ClientID,
	//		UsernameClaim:  additionalOidcConfig.UsernameClaim,
	//		UsernamePrefix: additionalOidcConfig.UsernamePrefix,
	//		GroupsClaim:    additionalOidcConfig.GroupsClaim,
	//		GroupsPrefix:   additionalOidcConfig.GroupsPrefix,
	//		RequiredClaims: additionalOidcConfig.RequiredClaims,
	//	},
	//}

	cr := &authenticationv1alpha1.OpenIDConnect{
		TypeMeta: metav1.TypeMeta{
			Kind:       "OpenIDConnect",
			APIVersion: "authentication.gardener.cloud/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-oidc",
		},
		Spec: authenticationv1alpha1.OIDCAuthenticationSpec{
			IssuerURL:      "https://token.actions.githubusercontent.com",
			ClientID:       "my-kubernetes-cluster",
			UsernameClaim:  ptr.To("sub"),
			UsernamePrefix: ptr.To("actions-oidc:"),
			GroupsPrefix:   ptr.To("deploy-kubernetes"),
			GroupsClaim:    ptr.To("myOrg/myRepo"),
			RequiredClaims: map[string]string{},
		},
	}

	return cr
}
