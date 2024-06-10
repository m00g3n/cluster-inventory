package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

func ExtendWithAnnotationsAndLabels(runtime imv1.Runtime, shoot *gardener.Shoot) error {
	shoot.Labels = getLabels(runtime)
	shoot.Annotations = getAnnotations(runtime)

	return nil
}

func getAnnotations(runtime imv1.Runtime) map[string]string {
	annotations := map[string]string{
		"infrastructuremanager.kyma-project.io/runtime-id": runtime.Labels["kyma-project.io/runtime-id:"],
	}

	if runtime.Spec.Shoot.LicenceType != nil && *runtime.Spec.Shoot.LicenceType != "" {
		annotations["infrastructuremanager.kyma-project.io/licence-type"] = *runtime.Spec.Shoot.LicenceType
	}

	return annotations
}

func getLabels(runtime imv1.Runtime) map[string]string {
	labels := map[string]string{
		"account":    runtime.Labels["kyma-project.io/global-account-id"],
		"subaccount": runtime.Labels["kyma-project.io/subaccount-id"],
	}

	return labels
}
