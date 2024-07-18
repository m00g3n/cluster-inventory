package fsm

import (
	"context"
	runtime_comparer "github.com/kyma-project/infrastructure-manager/internal/testing/runtime-comparer"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnCreateShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Create shoot")

	if m.SaveToPV {
		err := runtime_comparer.WriteToPV(s.instance, m.PVCPath)
		if err != nil {
			m.log.Error(err, "Failed to write Runtime CR to Persistent Volume")

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonProcessingErr,
				"False",
				"Failed to write Runtime CR to Persistent Volume",
			)
			return updateStatusAndRequeueAfter(gardenerRequeueDuration)
		}
	}

	newShoot, err := convertShoot(&s.instance, m.ConverterConfig)
	if err != nil {
		m.log.Error(err, "Failed to convert Runtime instance to shoot object")
		return updateStatePendingWithErrorAndStop(&s.instance, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonConversionError, "Runtime conversion error")
	}

	err = m.ShootClient.Create(ctx, &newShoot)

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

	s.instance.UpdateStatePending(
		imv1.ConditionTypeRuntimeProvisioned,
		imv1.ConditionReasonShootCreationPending,
		"Unknown",
		"Shoot is pending",
	)

	shouldPersistShoot := m.PVCPath != ""
	if shouldPersistShoot {
		s.shoot = newShoot.DeepCopy()
		return switchState(sFnPersistShoot)
	}

	return updateStatusAndRequeueAfter(gardenerRequeueDuration)
}
