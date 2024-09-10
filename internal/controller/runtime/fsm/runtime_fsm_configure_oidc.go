package fsm

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureOidc(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure OIDC state")

	//TODO: 0. plug in to the state machine without actually doing anything
	//TODO: 2/3. check condition if going into this state is necessary
	//TODO: 2/3. remove existing "additionalOidcs" from the Shoot spec
	//TODO: 3.5 revisit the default OIDC requirements, implementation drafted in https://github.com/kyma-project/infrastructure-manager/commit/32f016d38dafdee0ff1f9e772183825c703f046d
	//TODO: 1. recreate "additionalOidcs" in the Shoot spec
	//TODO: 4. create OpenID resource on the shoot
	//TODO: 4. extract rest client to a separate interface
	//TODO: 5. unit test the function
	// create OIDCConnect CR in shoot

	//s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonGardenerCRCreated, "Unknown", "Gardener Cluster CR created, waiting for readiness")
	//return updateStatusAndRequeueAfter(controlPlaneRequeueDuration)

	// wait section - is it needed?

	s.instance.UpdateStateReady(
		imv1.ConditionTypeOidcConfigured,
		imv1.ConditionReasonConfigurationCompleted,
		"OIDC configuration completed",
	)

	m.log.Info("OIDC has been configured", "Name", s.shoot.Name)

	return updateStatusAndStop()
}
