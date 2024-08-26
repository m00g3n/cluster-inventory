package fsm

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureOidc(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure OIDC state")

	//TODO: 2/3.check if doing this is necessary
	//TODO: 2/3. remove existing "additionalOidcs" from the Shoot spec
	//TODO: 1. recreate "additionalOidcs" in the Shoot spec
	//TODO: 4. extract rest client to a separate interface
	//TODO: 5. unit test the function
	// create OIDCConnect CR in shoot

	//s.instance.UpdateStatePending(imv1.ConditionTypeRuntimeKubeconfigReady, imv1.ConditionReasonGardenerCRCreated, "Unknown", "Gardener Cluster CR created, waiting for readiness")
	//return updateStatusAndRequeueAfter(controlPlaneRequeueDuration)

	// wait section - is it needed?

	//m.log.Info("OIDC has been configured", "Name", runtimeID)

	return ensureStatusConditionIsSetAndContinue(&s.instance,
		imv1.ConditionTypeRuntimeKubeconfigReady,
		imv1.ConditionReasonGardenerCRReady,
		"Gardener Cluster CR is ready.",
		sFnApplyClusterRoleBindings)
}
