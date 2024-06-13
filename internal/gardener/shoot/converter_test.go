package shoot

import (
	"testing"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot/hyperscaler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		Kubernetes: KubernetesConfig{
			DefaultVersion: "1.29",
		},
		DNS: DNSConfig{
			SecretName:   "dns-secret",
			DomainPrefix: "dev.mydomain.com",
			ProviderType: "aws-route53",
		},
		Provider: ProviderConfig{
			AWS: AWSConfig{
				EnableIMDSv2: true,
			},
		},
	}
}

func fixRuntime() imv1.Runtime {
	kubernetesVersion := "1.28"
	clientID := "client-id"
	groupsClaim := "groups"
	issuerURL := "https://my.cool.tokens.com"
	usernameClaim := "sub"

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
					Type: hyperscaler.TypeAWS,
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
						OidcConfig: gardener.OIDCConfig{
							ClientID:    &clientID,
							GroupsClaim: &groupsClaim,
							IssuerURL:   &issuerURL,
							SigningAlgs: []string{
								"RS256",
							},
							UsernameClaim: &usernameClaim,
						},
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
