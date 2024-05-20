package extenders

import (
	"testing"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOidcExtender(t *testing.T) {
	t.Run("Create kubernetes config", func(t *testing.T) {
		// given
		kubernetesVersion := "1.28"
		clientID := "client-id"
		groupsClaim := "groups"
		issuerURL := "https://my.cool.tokens.com"
		usernameClaim := "sub"

		shoot := fixEmptyGardenerShoot("test", "kcp-system")
		runtimeShoot := imv1.RuntimeShoot{
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
		}

		// when
		extender := NewExtendWithKubernetes("")
		err := extender(runtimeShoot, &shoot)

		// then
		require.NoError(t, err)
		assert.Equal(t, *runtimeShoot.Kubernetes.Version, shoot.Spec.Kubernetes.Version)
		assert.Equal(t, runtimeShoot.Kubernetes.KubeAPIServer.OidcConfig, *shoot.Spec.Kubernetes.KubeAPIServer.OIDCConfig)
	})
}
