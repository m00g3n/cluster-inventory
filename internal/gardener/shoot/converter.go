package shoot

import (
	"fmt"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/extender"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Extend func(imv1.Runtime, *gardener.Shoot) error

type Converter struct {
	extenders []Extend
	config    ConverterConfig
}

type ProviderConfig struct {
	AWS AWSConfig
}

type AWSConfig struct {
	EnableIMDSv2 bool
}

type DNSConfig struct {
	SecretName   string
	DomainPrefix string
	ProviderType string
}

type KubernetesConfig struct {
	DefaultVersion string
}

type ConverterConfig struct {
	Kubernetes   KubernetesConfig
	DNS          DNSConfig
	Provider     ProviderConfig
	MachineImage MachineImageConfig
	Gardener     GardenerConfig
}

type GardenerConfig struct {
	ProjectName string
}

type MachineImageConfig struct {
	DefaultVersion string
}

func NewConverter(config ConverterConfig) Converter {
	extenders := []Extend{
		extender.ExtendWithAnnotations,
		extender.ExtendWithLabels,
		extender.NewKubernetesVersionExtender(config.Kubernetes.DefaultVersion),
		extender.NewProviderExtender(config.Provider.AWS.EnableIMDSv2, config.MachineImage.DefaultVersion),
		extender.NewDNSExtender(config.DNS.SecretName, config.DNS.DomainPrefix, config.DNS.ProviderType),
		extender.ExtendWithOIDC,
		extender.ExtendWithCloudProfile,
		extender.ExtendWithNetworkFilter,
		extender.ExtendWithExposureClassName,
	}

	return Converter{
		extenders: extenders,
		config:    config,
	}
}

func (c Converter) ToShoot(runtime imv1.Runtime) (gardener.Shoot, error) {
	// The original implementation in the Provisioner: https://github.com/kyma-project/control-plane/blob/3dd257826747384479986d5d79eb20f847741aa6/components/provisioner/internal/model/gardener_config.go#L127

	// If you need to enhance the converter please adhere to the following convention:
	// - fields taken directly from Runtime CR must be added in this function
	// - if any logic is needed to be implemented, either enhance existing, or create a new extender

	shoot := gardener.Shoot{
		ObjectMeta: v1.ObjectMeta{
			Name:      runtime.Spec.Shoot.Name,
			Namespace: fmt.Sprintf("garden-%s", c.config.Gardener.ProjectName),
		},
		Spec: gardener.ShootSpec{
			Purpose:           &runtime.Spec.Shoot.Purpose,
			Region:            runtime.Spec.Shoot.Region,
			SecretBindingName: &runtime.Spec.Shoot.SecretBindingName,
			Networking: &gardener.Networking{
				Type:     runtime.Spec.Shoot.Networking.Type,
				Nodes:    &runtime.Spec.Shoot.Networking.Nodes,
				Pods:     &runtime.Spec.Shoot.Networking.Pods,
				Services: &runtime.Spec.Shoot.Networking.Services,
			},
			ControlPlane: &runtime.Spec.Shoot.ControlPlane,
		},
	}

	for _, extend := range c.extenders {
		if err := extend(runtime, &shoot); err != nil {
			return gardener.Shoot{}, err
		}
	}

	return shoot, nil
}
