package fsm

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnProcessShoot(_ context.Context, _ *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	// process shoot get kubeconfig and create cluster role bindings
	s.instance.UpdateStateReady(
		imv1.ConditionTypeRuntimeProvisioned,
		imv1.ConditionReasonConfigurationCompleted,
		"Runtime creation completed successfully")

	return stopWithNoRequeue()
}
