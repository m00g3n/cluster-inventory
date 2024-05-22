package azure

import (
	"encoding/json"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const infrastructureConfigKind = "InfrastructureConfig"
const controlPlaneConfigKind = "ControlPlaneConfig"
const apiVersion = "azure.provider.extensions.gardener.cloud/v1alpha1"

func GetInfrastructureConfig(workerCIDR string, zones []string) ([]byte, error) {
	return json.Marshal(NewInfrastructureConfig(workerCIDR, zones))
}

func GetControlPlaneConfig() ([]byte, error) {
	return json.Marshal(NewControlPlaneConfig())
}

func NewControlPlaneConfig() *ControlPlaneConfig {
	return &ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       controlPlaneConfigKind,
			APIVersion: apiVersion,
		},
	}
}

func NewInfrastructureConfig(workerCIDR string, zones []string) InfrastructureConfig {
	isZoned := len(zones) > 0

	azureConfig := InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       infrastructureConfigKind,
			APIVersion: apiVersion,
		},
		Networks: NetworkConfig{
			VNet: VNet{
				CIDR: &workerCIDR,
			},
		},
		Zoned: isZoned,
	}

	azureConfig.Networks.Zones = generateAzureZones(workerCIDR, zones)

	return azureConfig
}
