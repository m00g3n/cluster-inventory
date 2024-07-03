package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnCreateShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if s.shoot == nil {
		m.log.Info("Gardener shoot does not exist, creating new one")
		newShoot, err := convertShoot(&s.instance)
		if err != nil {
			m.log.Error(err, "Failed to convert Runtime instance to shoot object")
			return updateStatePendingWithErrorAndStop(&s.instance, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonConversionError, "Runtime conversion error")
		}

		_, err = m.ShootClient.Create(ctx, &newShoot, v1.CreateOptions{})

		if err != nil {
			m.log.Error(err, "Failed to create new gardener Shoot")

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonGardenerError,
				"False",
				"Gardener API create error",
			)
			return updateStatusAndRequeueAfter(gardenerRequeueDuration)
		}
		m.log.Info("Gardener shoot for runtime initialised successfully", "Name", newShoot.Name, "Namespace", newShoot.Namespace)
	}

	s.instance.UpdateStatePending(
		imv1.ConditionTypeRuntimeProvisioned,
		imv1.ConditionReasonShootCreationPending,
		"Unknown",
		"Shoot is pending",
	)

	shouldPersistShoot := m.PVCPath != ""
	if shouldPersistShoot {
		return switchState(sFnPersistShoot)
	}

	return updateStatusAndRequeueAfter(gardenerRequeueDuration)
}
