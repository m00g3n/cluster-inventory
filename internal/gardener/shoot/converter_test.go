package shoot

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestConverter(t *testing.T) {
	t.Run("Create shoot from Runtime", func(t *testing.T) {
		// given
		runtime := fixRuntime()
		converterConfig := fixConverterConfig()
		converter := NewConverter(converterConfig)

		// when
		shoot, err := converter.ToShoot(runtime)

		// then
		require.NoError(t, err)
		assert.Equal(t, runtime.Spec.Shoot.Purpose, *shoot.Spec.Purpose)
		assert.Equal(t, runtime.Spec.Shoot.Region, shoot.Spec.Region)
		assert.Equal(t, runtime.Spec.Shoot.SecretBindingName, *shoot.Spec.SecretBindingName)
		assert.Equal(t, runtime.Spec.Shoot.ControlPlane, *shoot.Spec.ControlPlane)
	})
}

func fixConverterConfig() ConverterConfig {
	return ConverterConfig{
		DefaultKubernetesVersion: "1.29",
		DNSSecretName:            "dns-secret",
		DomainPrefix:             "dev.mydomain.com",
	}
}

func fixRuntime() imv1.Runtime {
	kubernetesVersion := "1.28"

	return imv1.Runtime{
		ObjectMeta: v1.ObjectMeta{
			Name:      "runtime",
			Namespace: "kcp-system",
		},
		Spec: imv1.RuntimeSpec{
			Shoot: imv1.RuntimeShoot{
				Purpose:           "production",
				Region:            "eu-central-1",
				SecretBindingName: "my-secret",
				Provider: imv1.Provider{
					Type: "aws",
					Workers: []gardener.Worker{
						{
							Name: "worker",
							Machine: gardener.Machine{
								Type: "m6i.large",
							},
							Minimum: 1,
							Maximum: 3,
							Zones: []string{
								"eu-central-1a",
								"eu-central-1b",
								"eu-central-1c",
							},
						},
					},
				},
				Kubernetes: imv1.Kubernetes{
					Version: &kubernetesVersion,
					KubeAPIServer: imv1.APIServer{
						OidcConfig: fixGardenerOidcConfig(),
					},
				},
				Networking: imv1.Networking{
					Pods:     "100.64.0.0/12",
					Nodes:    "10.250.0.0/16",
					Services: "100.104.0.0/13",
				},
				ControlPlane: gardener.ControlPlane{
					HighAvailability: &gardener.HighAvailability{
						FailureTolerance: gardener.FailureTolerance{
							Type: gardener.FailureToleranceTypeZone,
						},
					},
				},
			},
		},
	}
}

func fixGardenerOidcConfig() gardener.OIDCConfig {
	clientID := "client-id"
	groupsClaim := "groups"
	issuerURL := "https://my.cool.tokens.com"
	usernameClaim := "sub"

	return gardener.OIDCConfig{
		ClientID:    &clientID,
		GroupsClaim: &groupsClaim,
		IssuerURL:   &issuerURL,
		SigningAlgs: []string{
			"RS256",
		},
		UsernameClaim: &usernameClaim,
	}
}

func fixEmptyGardenerShoot(name, namespace string) gardener.Shoot {
	return gardener.Shoot{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: gardener.ShootSpec{},
	}
}
