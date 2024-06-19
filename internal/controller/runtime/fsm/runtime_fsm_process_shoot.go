package fsm

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnProcessShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	if !s.instance.IsRuntimeStateSet(imv1.RuntimeStatePending, imv1.ConditionTypeRuntimeConfigured, imv1.ConditionReasonConfigurationStarted) {
		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeConfigured,
			imv1.ConditionReasonConfigurationStarted,
			"Unknown",
			"Runtime processing initialised",
		)

		return stopWithRequeue()
	}

	// TODO: now let's process shoot get kubeconfig and create cluster role bindings

	s.instance.UpdateStateReady(
		imv1.ConditionTypeRuntimeProvisioned,
		imv1.ConditionReasonConfigurationCompleted,
		"Runtime creation completed successfully")

	return stopWithNoRequeue()
}
