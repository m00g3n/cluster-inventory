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
	Kubernetes KubernetesConfig
	DNS        DNSConfig
	Provider   ProviderConfig
}

func NewConverter(config ConverterConfig) Converter {
	extenders := []Extend{
		extender.ExtendWithAnnotations,
		extender.ExtendWithLabels,
		extender.NewKubernetesVersionExtender(config.Kubernetes.DefaultVersion),
		extender.ExtendWithNetworking,
		extender.NewProviderExtender(config.Provider.AWS.EnableIMDSv2),
		extender.NewExtendWithDNS(config.DNS.SecretName, config.DNS.DomainPrefix, config.DNS.ProviderType),
		extender.ExtendWithOIDC,
	}

	return Converter{
		extenders: extenders,
	}
}

func (c Converter) ToShoot(runtime imv1.Runtime) (gardener.Shoot, error) {
	// The original implementation in the Provisioner: https://github.com/kyma-project/control-plane/blob/3dd257826747384479986d5d79eb20f847741aa6/components/provisioner/internal/model/gardener_config.go#L127
	// Note: shoot.Spec.ExposureClassNames field is ignored as KEB didn't send this field to the Provisioner

	shoot := gardener.Shoot{
		ObjectMeta: v1.ObjectMeta{
			Name:      runtime.Spec.Shoot.Name,
			Namespace: fmt.Sprintf("garden-%s", "kyma-dev"), //nolint:godox TODO: make it more dynamic - this should be the gardener project namespace
		},
		Spec: gardener.ShootSpec{
			Purpose:           &runtime.Spec.Shoot.Purpose,
			Region:            runtime.Spec.Shoot.Region,
			SecretBindingName: &runtime.Spec.Shoot.SecretBindingName,
			ControlPlane:      &runtime.Spec.Shoot.ControlPlane, //nolint:godox TODO: check HAavailability (also in migrator)
		},
	}

	for _, extend := range c.extenders {
		if err := extend(runtime, &shoot); err != nil {
			return gardener.Shoot{}, err
		}
	}

	return shoot, nil
}

func PtrTo[T any](v T) *T {
	return &v
}
