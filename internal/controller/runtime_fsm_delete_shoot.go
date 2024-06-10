package controller

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/apimachinery/pkg/api/meta"

	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnDeleteShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if !isDeleting(s) {
		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeDeprovisioning,
			imv1.ConditionReasonDeletion,
			"Runtime deletion initialised",
		)

		return stopWithRequeue()
	}

	return stopWithNoRequeue()
}

func isDeleting(s *systemState) bool {
	condition := meta.FindStatusCondition(s.instance.Status.Conditions, string(imv1.ConditionTypeRuntimeDeprovisioning))
	if condition == nil {
		return false
	}

	if condition.Reason != string(imv1.ConditionReasonDeletion) &&
		condition.Reason != string(imv1.ConditionReasonDeletionErr) {
		return false
	}

	return true
}
