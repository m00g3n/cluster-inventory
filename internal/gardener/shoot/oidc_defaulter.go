package shoot

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

func DefaultOidcIfNotPresent(runtime *imv1.Runtime, oidcProviderCfg OidcProvider) {
	oidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig

	if ShouldDefaultOidcConfig(oidcConfig) {
		defaultOIDCConfig := CreateDefaultOIDCConfig(oidcProviderCfg)
		runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig = defaultOIDCConfig
	}
}

func ShouldDefaultOidcConfig(config gardener.OIDCConfig) bool {
	return config.ClientID == nil && config.IssuerURL == nil
}

func CreateDefaultOIDCConfig(defaultSharedIASTenant OidcProvider) gardener.OIDCConfig {
	return gardener.OIDCConfig{
		ClientID:       &defaultSharedIASTenant.ClientID,
		GroupsClaim:    &defaultSharedIASTenant.GroupsClaim,
		IssuerURL:      &defaultSharedIASTenant.IssuerURL,
		SigningAlgs:    defaultSharedIASTenant.SigningAlgs,
		UsernameClaim:  &defaultSharedIASTenant.UsernameClaim,
		UsernamePrefix: &defaultSharedIASTenant.UsernamePrefix,
	}
}
