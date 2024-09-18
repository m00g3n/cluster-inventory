package shoot

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

func DefaultAdditionalOidcIfNotPresent(runtime *imv1.Runtime) {
	oidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig
	additionalOidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig

	if nil == additionalOidcConfig {
		additionalOidcConfig = &[]gardener.OIDCConfig{}
		*additionalOidcConfig = append(*additionalOidcConfig, oidcConfig)
	}

	runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig = additionalOidcConfig
}
