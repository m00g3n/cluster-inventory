package controller

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnDeleteShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if !s.instance.IsRuntimeStateSet(imv1.RuntimeStateDeleting, imv1.ConditionTypeRuntimeDeprovisioning, imv1.ConditionReasonDeletion) {
		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeDeprovisioning,
			imv1.ConditionReasonDeletion,
			"Runtime deletion initialised",
		)
		return stopWithRequeue()
	}

	return stopWithNoRequeue()
}
