package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

const (
	DefaultAWSEuAccessAWSRegion = "eu-central-1"
	DefaultEuAccessAzureRegion  = "switzerlandnorth"
)

func ExtendWithRegion(runtime imv1.Runtime, shoot *gardener.Shoot) error {
	shoot.Spec.Region = getRegion(runtime)

	return nil
}

func getRegion(runtime imv1.Runtime) string {
	if isEuAccess(runtime.Spec.Shoot.PlatformRegion) {
		switch runtime.Spec.Shoot.Provider.Type {
		case ProviderTypeAWS:
			{
				return DefaultAWSEuAccessAWSRegion
			}
		case ProviderTypeAzure:
			{
				return DefaultEuAccessAzureRegion
			}
		}
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
