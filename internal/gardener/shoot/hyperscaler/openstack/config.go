package openstack

import (
	"encoding/json"

	"github.com/gardener/gardener-extension-provider-openstack/pkg/apis/openstack/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	infrastructureConfigKind    = "InfrastructureConfig"
	controlPlaneConfigKind      = "ControlPlaneConfig"
	apiVersion                  = "openstack.provider.extensions.gardener.cloud/v1alpha1"
	defaultFloatingPoolName     = "FloatingIP-external-kyma-01"
	defaultLoadBalancerProvider = "f5"
)

func GetInfrastructureConfig(workerCIDR string, _ []string) ([]byte, error) {
	return json.Marshal(NewInfrastructureConfig(workerCIDR))
}

func GetControlPlaneConfig(_ []string) ([]byte, error) {
	return json.Marshal(NewControlPlaneConfig())
}

func NewInfrastructureConfig(workerCIDR string) v1alpha1.InfrastructureConfig {
	return v1alpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       infrastructureConfigKind,
			APIVersion: apiVersion,
		},
		FloatingPoolName: defaultFloatingPoolName,
		Networks: v1alpha1.Networks{
			Workers: workerCIDR,
		},
	}
}

func NewControlPlaneConfig() *v1alpha1.ControlPlaneConfig {
	return &v1alpha1.ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       controlPlaneConfigKind,
			APIVersion: apiVersion,
		},
		LoadBalancerProvider: defaultLoadBalancerProvider,
	}
}
