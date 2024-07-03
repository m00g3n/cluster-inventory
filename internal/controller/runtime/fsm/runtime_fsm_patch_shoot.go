package fsm

import (
	"context"
	"encoding/json"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	gardener_shoot "github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnPatchExistingShoot(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Patch shoot state")

	updatedShoot, err := convertShoot(&s.instance)
	if err != nil {
		m.log.Error(err, "Failed to convert Runtime instance to shoot object, exiting with no retry")
		return updateStatePendingWithErrorAndStop(&s.instance, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonConversionError, "Runtime conversion error")
	}
	m.log.Info("Shoot converted successfully", "Name", updatedShoot.Name, "Namespace", updatedShoot.Namespace)

	shootData, err := json.Marshal(updatedShoot)
	if err != nil {
		m.log.Error(err, "Failed to marshal shoot object to JSON, exiting with no retry")
		return updateStatePendingWithErrorAndStop(&s.instance, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonConversionError, "Shoot marshaling error")
	}
	m.log.Info("Shoot marshaled successfully", "Name", updatedShoot.Name, "Namespace", updatedShoot.Namespace)

	patched, err := m.ShootClient.Patch(context.Background(), s.shoot.Name, types.ApplyPatchType, shootData, v1.PatchOptions{FieldManager: "kim", Force: ptrTo(true)})

	if err != nil {
		m.log.Error(err, "Failed to patch shoot object with JSON, exiting with no retry")
		return updateStatePendingWithErrorAndStop(&s.instance, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonGardenerError, "Shoot patch error")
	}

	if patched.Generation == s.shoot.Generation {
		m.log.Info("Gardener shoot for runtime did not change after patch, moving to processing", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)
		return switchState(sFnProcessShoot)
	} else {
		m.log.Info("Gardener shoot for runtime patched successfully", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)

		s.shoot = patched.DeepCopy()

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonProcessing,
			"Unknown",
			"Shoot is pending",
		)

		shouldPersistShoot := m.PVCPath != ""
		if shouldPersistShoot {
			return switchState(sFnPersistShoot)
		}

		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}
}

func convertShoot(instance *imv1.Runtime) (gardener.Shoot, error) {
	converterConfig := FixConverterConfig()
	converter := gardener_shoot.NewConverter(converterConfig)
	shoot, err := converter.ToShoot(*instance) // returned error is always nil BTW

	if err == nil {
		setObjectFields(&shoot)
	}

	return shoot, err
}

func FixConverterConfig() gardener_shoot.ConverterConfig {
	return gardener_shoot.ConverterConfig{
		Kubernetes: gardener_shoot.KubernetesConfig{
			DefaultVersion: "1.29", //nolint:godox TODO: Should be parametrised
		},

		DNS: gardener_shoot.DNSConfig{
			SecretName:   "aws-route53-secret-dev",
			DomainPrefix: "dev.kyma.ondemand.com",
			ProviderType: "aws-route53",
		},
		Provider: gardener_shoot.ProviderConfig{
			AWS: gardener_shoot.AWSConfig{
				EnableIMDSv2: true, //nolint:godox TODO: Should be parametrised
			},
		},
		Gardener: gardener_shoot.GardenerConfig{
			ProjectName: "kyma-dev", //nolint:godox TODO: should be parametrised
		},
	}
}

// workaround
func setObjectFields(shoot *gardener.Shoot) {
	shoot.Kind = "Shoot"
	shoot.APIVersion = "core.gardener.cloud/v1beta1"
	shoot.ManagedFields = nil
}

func updateStatePendingWithErrorAndStop(instance *imv1.Runtime, c imv1.RuntimeConditionType, r imv1.RuntimeConditionReason, msg string) (stateFn, *ctrl.Result, error) {
	instance.UpdateStatePending(c, r, "False", msg)
	return updateStatusAndStop()
}

func ptrTo[T any](v T) *T {
	return &v
}
