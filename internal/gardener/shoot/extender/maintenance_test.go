package extender

import (
	"testing"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestMaintenanceExtender(t *testing.T) {
	for _, testCase := range []struct {
		name                string
		kubernetesVersion   bool
		machineImageVersion bool
	}{
		{
			name:                "Enable auto-update for only kubernetesVersion",
			kubernetesVersion:   true,
			machineImageVersion: false,
		},
		{
			name:                "Enable auto-update for only machineImageVersion",
			kubernetesVersion:   false,
			machineImageVersion: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			shoot := fixEmptyGardenerShoot("test", "dev")
			shoot.Spec.Maintenance = &gardener.Maintenance{
				AutoUpdate: &gardener.MaintenanceAutoUpdate{},
			}
			runtimeShoot := imv1.Runtime{
				Spec: imv1.RuntimeSpec{
					Shoot: imv1.RuntimeShoot{
						Name: "test",
					},
				},
			}

			// when
			extender := NewMaintenanceExtender(testCase.kubernetesVersion, testCase.machineImageVersion)
			err := extender(runtimeShoot, &shoot)

			// then
			assert.NoError(t, err)
			assert.Equal(t, testCase.kubernetesVersion, shoot.Spec.Maintenance.AutoUpdate.KubernetesVersion)
			assert.Equal(t, testCase.machineImageVersion, *shoot.Spec.Maintenance.AutoUpdate.MachineImageVersion)
		})
	}
}
