package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

const DefaultEuAccessAWSRegion = "eu-central-1"

func ExtendWithRegion(runtime imv1.Runtime, shoot *gardener.Shoot) error {
	shoot.Spec.Region = getRegion(runtime)

	return nil
}

func getRegion(runtime imv1.Runtime) string {
	if isEuAccess(runtime.Spec.Shoot.PlatformRegion) {
		return DefaultEuAccessAWSRegion
	}

	return runtime.Spec.Shoot.Region
}

func isEuAccess(platformRegion string) bool {
	switch platformRegion {
	case "cf-eu11":
		return true
	case "cf-ch20":
		return true
	}
	return false
}
