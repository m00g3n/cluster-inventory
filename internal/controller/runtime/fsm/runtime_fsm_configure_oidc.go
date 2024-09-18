package fsm

import (
	"context"
	"fmt"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	authenticationv1alpha1 "github.com/gardener/oidc-webhook-authenticator/apis/authentication/v1alpha1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/extender"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureOidc(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure OIDC state")

	if !isOidcExtensionEnabled(*s.shoot) {
		m.log.Info("OIDC extension is disabled")
		s.instance.UpdateStateReady(
			imv1.ConditionTypeOidcConfigured,
			imv1.ConditionReasonOidcConfigured,
			"OIDC extension disabled",
		)
		return updateStatusAndStop()
	}
	shoot.DefaultAdditionalOidcIfNotPresent(&s.instance)
	validationError := validateOidcConfiguration(s.instance)
	if validationError != nil {
		m.log.Error(validationError, "default OIDC configuration is not present")
		updateConditionFailed(&s.instance)
		return updateStatusAndStopWithError(validationError)
	}

	err := createOpenIDConnectResources(ctx, m, s)

	if k8serrors.IsAlreadyExists(err) {
		return switchState(sFnApplyClusterRoleBindings)
	}

	if err != nil {
		m.log.Error(err, "Failed to create OpenIDConnect resource")
		updateConditionFailed(&s.instance)
		return updateStatusAndStopWithError(err)
	}

	s.instance.UpdateStateReady(
		imv1.ConditionTypeOidcConfigured,
		imv1.ConditionReasonOidcConfigured,
		"OIDC configuration completed",
	)

	m.log.Info("OIDC has been configured", "Name", s.shoot.Name)

	return switchState(sFnApplyClusterRoleBindings)
}

func validateOidcConfiguration(rt imv1.Runtime) (err error) {
	if rt.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig == nil {
		err = errors.New("default OIDC configuration is not present")
	}
	return err
}

func createOpenIDConnectResources(ctx context.Context, m *fsm, s *systemState) error {
	srscClient := m.ShootClient.SubResource("adminkubeconfig")
	shootAdminClient, shootClientError := GetShootClient(ctx, srscClient, s.shoot)

	if shootClientError != nil {
		return shootClientError
	}

	additionalOidcConfigs := *s.instance.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig

	var errResourceCreation error
	for id, additionalOidcConfig := range additionalOidcConfigs {
		openIDConnectResource := createOpenIDConnectResource(additionalOidcConfig, id)
		errResourceCreation = shootAdminClient.Create(ctx, openIDConnectResource)
	}
	return errResourceCreation
}

func isOidcExtensionEnabled(shoot gardener.Shoot) bool {
	for _, extension := range shoot.Spec.Extensions {
		if extension.Type == extender.OidcExtensionType {
			return !(*extension.Disabled)
		}
	}
	return false
}

func createOpenIDConnectResource(additionalOidcConfig gardener.OIDCConfig, oidcID int) *authenticationv1alpha1.OpenIDConnect {
	cr := &authenticationv1alpha1.OpenIDConnect{
		TypeMeta: metav1.TypeMeta{
			Kind:       "OpenIDConnect",
			APIVersion: "authentication.gardener.cloud/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("kyma-oidc-%v", oidcID),
		},
		Spec: authenticationv1alpha1.OIDCAuthenticationSpec{
			IssuerURL:      *additionalOidcConfig.IssuerURL,
			ClientID:       *additionalOidcConfig.ClientID,
			UsernameClaim:  additionalOidcConfig.UsernameClaim,
			UsernamePrefix: additionalOidcConfig.UsernamePrefix,
			GroupsClaim:    additionalOidcConfig.GroupsClaim,
			GroupsPrefix:   additionalOidcConfig.GroupsPrefix,
			RequiredClaims: additionalOidcConfig.RequiredClaims,
		},
	}

	return cr
}

func updateConditionFailed(rt *imv1.Runtime) {
	rt.UpdateStatePending(
		imv1.ConditionTypeOidcConfigured,
		imv1.ConditionReasonOidcError,
		string(metav1.ConditionFalse),
		"failed to configure OIDC",
	)
}
