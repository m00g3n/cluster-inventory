package azure

import (
	"encoding/json"

	"github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const infrastructureConfigKind = "InfrastructureConfig"
const controlPlaneConfigKind = "ControlPlaneConfig"
const apiVersion = "azure.provider.extensions.gardener.cloud/v1alpha1"

func GetInfrastructureConfig(workerCIDR string, zones []string) ([]byte, error) {
	return json.Marshal(NewInfrastructureConfig(workerCIDR, zones))
}

func GetControlPlaneConfig(_ []string) ([]byte, error) {
	return json.Marshal(NewControlPlaneConfig())
}

func NewControlPlaneConfig() *v1alpha1.ControlPlaneConfig {
	return &v1alpha1.ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       controlPlaneConfigKind,
			APIVersion: apiVersion,
		},
	}
}

func NewInfrastructureConfig(workerCIDR string, zones []string) v1alpha1.InfrastructureConfig {
	// All Azure shoots are zoned.
	// No zones - the shoot configuration is invalid.
	// We should validate the config before calling this function.
	isZoned := len(zones) > 0

	azureConfig := v1alpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       infrastructureConfigKind,
			APIVersion: apiVersion,
		},
		Networks: v1alpha1.NetworkConfig{
			VNet: v1alpha1.VNet{
				CIDR: &workerCIDR,
			},
		},
		Zoned: isZoned,
	}

	azureConfig.Networks.Zones = generateAzureZones(workerCIDR, zones)

	return azureConfig
}
