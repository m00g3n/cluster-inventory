package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

const NetworkFilterType = "shoot-networking-filter"

func ExtendWithNetworkFilter(runtime imv1.Runtime, shoot *gardener.Shoot) error {
	networkingFilter := gardener.Extension{
		Type:     NetworkFilterType,
		Disabled: ToPtr(false),
	}

	shoot.Spec.Extensions = append(shoot.Spec.Extensions, networkingFilter)

	return nil
}

func ToPtr[T any](v T) *T {
	return &v
}
