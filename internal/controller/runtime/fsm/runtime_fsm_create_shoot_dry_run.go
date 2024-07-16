package fsm

import (
	"context"
	"fmt"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnCreateShootDryRun(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Create shoot [dry-run]")

	newShoot, err := convertShoot(&s.instance, m.ConverterConfig)
	if err != nil {
		m.log.Error(err, "Failed to convert Runtime instance to shoot object [dry-run]")
		return updateStatePendingWithErrorAndStop(
			&s.instance,
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonConversionError,
			"Runtime conversion error")
	}

	s.shoot = &newShoot
	s.instance.UpdateStateReady(
		imv1.ConditionTypeRuntimeProvisioned,
		imv1.ConditionReasonConfigurationCompleted,
		"Runtime processing completed successfully [dry-run]")

	// stop machine if persistence not enabled
	if m.PVCPath == "" {
		return updateStatusAndStop()
	}

	path := fmt.Sprintf("%s/%s-%s.yaml", m.PVCPath, s.shoot.Namespace, s.shoot.Name)
	if err := persist(path, s.shoot, m.writerProvider); err != nil {
		return updateStatusAndStopWithError(err)
	}
	return updateStatusAndStop()
}
