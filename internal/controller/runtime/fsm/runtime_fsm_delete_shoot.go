package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnDeleteShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if !s.instance.IsStateWithConditionSet(imv1.RuntimeStateTerminating, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonDeletion) {
		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonDeletion,
			"Unknown",
			"Runtime deletion initialised",
		)
		return stopWithRequeue()
	}

	err := m.ShootClient.Delete(ctx, s.instance.Name, v1.DeleteOptions{})

	if err != nil {
		m.log.Error(err, "Failed to delete gardener Shoot")

		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerError,
			"False",
			"Gardener API delete error",
		)
		return stopWithRequeue()
	}
	return stopWithNoRequeue()
}
