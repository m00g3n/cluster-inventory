package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/utils/ptr"
)

const (
	OidcExtensionType = "shoot-oidc-service"
)

// Extends Shoot spec with OIDC configuration and mutates Runtime spec with necessary OIDC defaults if missing
func ExtendWithOIDC(runtime *imv1.Runtime, shoot *gardener.Shoot) error {
	oidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig

	defaultAdditionalOidcIfNotPresent(runtime)
	setOIDCExtension(shoot)
	setKubeAPIServerOIDCConfig(shoot, oidcConfig)

	return nil
}

func defaultAdditionalOidcIfNotPresent(runtime *imv1.Runtime) {
	oidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig
	additionalOidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig

	if nil == additionalOidcConfig {
		additionalOidcConfig = &[]gardener.OIDCConfig{}
		*additionalOidcConfig = append(*additionalOidcConfig, oidcConfig)
	}

	runtime.Spec.Shoot.Kubernetes.KubeAPIServer.AdditionalOidcConfig = additionalOidcConfig
}

func setOIDCExtension(shoot *gardener.Shoot) {
	oidcService := gardener.Extension{
		Type:     OidcExtensionType,
		Disabled: ptr.To(false),
	}

	shoot.Spec.Extensions = append(shoot.Spec.Extensions, oidcService)
}

func setKubeAPIServerOIDCConfig(shoot *gardener.Shoot, oidcConfig gardener.OIDCConfig) {
	shoot.Spec.Kubernetes.KubeAPIServer = &gardener.KubeAPIServerConfig{
		OIDCConfig: &gardener.OIDCConfig{
			CABundle:       oidcConfig.CABundle,
			ClientID:       oidcConfig.ClientID,
			GroupsClaim:    oidcConfig.GroupsClaim,
			GroupsPrefix:   oidcConfig.GroupsPrefix,
			IssuerURL:      oidcConfig.IssuerURL,
			RequiredClaims: oidcConfig.RequiredClaims,
			SigningAlgs:    oidcConfig.SigningAlgs,
			UsernameClaim:  oidcConfig.UsernameClaim,
			UsernamePrefix: oidcConfig.UsernamePrefix,
		},
	}
}

//main OIDC task https://github.com/kyma-project/kyma/issues/18305#issuecomment-2128866460

func setOIDCConfig(shoot *gardener.Shoot, oidcConfig gardener.OIDCConfig) {
	shoot.Spec.Kubernetes.KubeAPIServer = &gardener.KubeAPIServerConfig{
		OIDCConfig: &oidcConfig,
	}
}

func getDefaultOIDCConfig() *gardener.OIDCConfig {
	return &gardener.OIDCConfig{

		// taken from https://github.tools.sap/kyma/management-plane-config/blob/20474fc793b147845b884160954d280f75b98a85/argoenv/keb/dev/values.yaml

		//TODO: move below's default configuration to:
		// - config file for local development
		// - management-plane-charts/config

		//CABundle:    //TODO: is it needed?
		ClientID:    ptr.To("xyz"),    //TODO: move to config file
		GroupsClaim: ptr.To("groups"), //TODO: move to config file
		//GroupsPrefix: TODO: is it needed?
		IssuerURL: ptr.To("https://kymatest.accounts400.ondemand.com"), //TODO: move to config file
		//RequiredClaims: TODO: is it needed?
		SigningAlgs:    []string{"RS256"}, //TODO: move to config file
		UsernameClaim:  ptr.To("sub"),     //TODO: move to config file
		UsernamePrefix: ptr.To("-"),       //TODO: move to config file
	}
}
