package shoot

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/extender"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Converter struct {
	extenders []extender.Extender
}

type ConverterConfig struct {
	DefaultKubernetesVersion string
	DNSSecretName            string
	DomainPrefix             string
	DNSProviderType          string
}

func NewConverter(config ConverterConfig) Converter {
	extenders := []extender.Extender{
		extender.AnnotationsExtender,
		extender.NewKubernetesExtender(config.DefaultKubernetesVersion),
		extender.NetworkingExtender,
		extender.ProviderExtender,
		extender.NewDNSExtender(config.DNSSecretName, config.DomainPrefix, config.DNSProviderType),
	}

	return Converter{
		extenders: extenders,
	}
}

func (c Converter) ToShoot(runtime imv1.Runtime) (gardener.Shoot, error) {
	shoot := gardener.Shoot{
		ObjectMeta: v1.ObjectMeta{
			Name:      runtime.Spec.Shoot.Name,
			Namespace: runtime.Namespace,
		},
		Spec: gardener.ShootSpec{
			Purpose:           &runtime.Spec.Shoot.Purpose,
			Region:            runtime.Spec.Shoot.Region,
			SecretBindingName: &runtime.Spec.Shoot.SecretBindingName,
			ControlPlane:      &runtime.Spec.Shoot.ControlPlane,
		},
	}

	for _, extend := range c.extenders {
		if err := extend(runtime.Spec.Shoot, &shoot); err != nil {
			return gardener.Shoot{}, err
		}
	}

	return shoot, nil
}
