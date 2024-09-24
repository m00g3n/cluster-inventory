package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/utils/ptr"
)

const (
	OidcExtensionType = "shoot-oidc-service"
)

func ShouldDefaultOidcConfig(config gardener.OIDCConfig) bool {
	return config.ClientID == nil && config.IssuerURL == nil
}

func NewOidcExtender(clientID, groupsClaim, issuerURL, usernameClaim, usernamePrefix string, signingAlgs []string) func(runtime imv1.Runtime, shoot *gardener.Shoot) error {
	return func(runtime imv1.Runtime, shoot *gardener.Shoot) error {
		if CanEnableExtension(runtime) {
			setOIDCExtension(shoot)
		}

		oidcConfig := runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig
		if ShouldDefaultOidcConfig(oidcConfig) {
			oidcConfig = gardener.OIDCConfig{
				ClientID:       &clientID,
				GroupsClaim:    &groupsClaim,
				IssuerURL:      &issuerURL,
				SigningAlgs:    signingAlgs,
				UsernameClaim:  &usernameClaim,
				UsernamePrefix: &usernamePrefix,
			}
		}
		setKubeAPIServerOIDCConfig(shoot, oidcConfig)

		return nil
	}
}
func CanEnableExtension(runtime imv1.Runtime) bool {
	canEnable := true
	createdByMigrator := runtime.Labels["operator.kyma-project.io/created-by-migrator"]

	if createdByMigrator == "true" {
		canEnable = false
	}

	return canEnable
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
