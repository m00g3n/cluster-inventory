package controller

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnProcessShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	//if !isProcessing(s) {
	//	s.instance.UpdateStateProcessing(
	//		imv1.ConditionTypeRuntimeProvisioning,
	//		imv1.ConditionReasonProcessing,
	//		"Runtime processing initialised",
	//	)
	//
	//	return stopWithRequeue()
	//}

	// TODO: now let's process shoot get kubeconfig and create cluster role bindings

	s.instance.UpdateStateProcessing(
		imv1.ConditionTypeRuntimeProvisioning,
		imv1.ConditionReasonProcessing,
		"Runtime processing completed successfully")

	return stopWithNoRequeue()
}
