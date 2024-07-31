package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/utils/ptr"
)

func NewMaintenanceExtender(kubernetesVersion, machineImageVersion bool) func(runtime imv1.Runtime, shoot *gardener.Shoot) error { //nolint:revive
	return func(runtime imv1.Runtime, shoot *gardener.Shoot) error {
		shoot.Spec.Maintenance = &gardener.Maintenance{
			AutoUpdate: &gardener.MaintenanceAutoUpdate{
				KubernetesVersion:   kubernetesVersion,
				MachineImageVersion: ptr.To(machineImageVersion),
			},
		}

		return nil
	}
}
