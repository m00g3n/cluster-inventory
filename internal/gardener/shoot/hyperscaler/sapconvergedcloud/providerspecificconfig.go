package sapconvergedcloud

import (
	"github.com/gardener/gardener-extension-provider-openstack/pkg/apis/openstack/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	infrastructureConfigKind = "InfrastructureConfig"
	controlPlaneConfigKind   = "ControlPlaneConfig"
	apiVersion               = "openstack.provider.extensions.gardener.cloud/v1alpha1"
)

func GetInfrastructureConfig(_ string, _ []string) ([]byte, error) {
	return nil, nil
}

func GetControlPlaneConfig(_ []string) ([]byte, error) {
	return nil, nil
}

func NewInfrastructureConfig(workerCIDR string) v1alpha1.InfrastructureConfig {
	return v1alpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       infrastructureConfigKind,
			APIVersion: apiVersion,
		},
	}
}

func NewControlPlaneConfig(zones []string) *v1alpha1.ControlPlaneConfig {
	return &v1alpha1.ControlPlaneConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       controlPlaneConfigKind,
			APIVersion: apiVersion,
		},
	}
}
