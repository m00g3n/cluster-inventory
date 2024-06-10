package controller

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	gardener_shoot "github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnCreateShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	converterConfig := fixConverterConfig()
	converter := gardener_shoot.NewConverter(converterConfig)
	shoot, err := converter.ToShoot(s.instance)

	if err != nil {
		m.log.Error(err, "unable to convert Runtime CR to a shoot object")

		s.instance.UpdateStateError(
			imv1.ConditionTypeRuntimeProvisioning,
			imv1.ConditionReasonConversionError,
			"Runtime conversion error",
		)

		return stopWithNoRequeue()
	}

	m.log.Info("Shoot mapped successfully", "Name", shoot.Name, "Namespace", shoot.Namespace, "Shoot", shoot)

	createdShoot, provisioningErr := m.shootClient.Create(ctx, &shoot, v1.CreateOptions{})

	if provisioningErr != nil {
		m.log.Error(provisioningErr, "Failed to create new gardener Shoot")

		s.instance.UpdateStateError(
			imv1.ConditionTypeRuntimeProvisioning,
			imv1.ConditionReasonCreationError,
			"Gardener API error",
		)

		return stopWithRequeue()
	}

	m.log.Info("Gardener shoot for runtime initialised successfully", "Name", createdShoot.Name, "Namespace", createdShoot.Namespace)

	s.instance.UpdateStateCreating(
		imv1.ConditionTypeRuntimeProvisioning,
		imv1.ConditionReasonShootCreationPending,
		"Shoot is pending",
	)

	return stopWithRequeue()
}
