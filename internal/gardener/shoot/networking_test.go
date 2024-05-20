package shoot

import (
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkingExtender(t *testing.T) {
	t.Run("Crete networking config", func(t *testing.T) {
		// given
		shoot := fixEmptyGardenerShoot("test", "kcp-system")
		runtimeShoot := imv1.RuntimeShoot{
			Networking: imv1.Networking{
				Pods:     "100.64.0.0/12",
				Nodes:    "10.250.0.0/16",
				Services: "100.104.0.0/13",
			},
		}

		// when
		err := networkingExtender(runtimeShoot, &shoot)

		// then
		require.NoError(t, err)
		assert.Equal(t, runtimeShoot.Networking.Nodes, *shoot.Spec.Networking.Nodes)
		assert.Equal(t, runtimeShoot.Networking.Pods, *shoot.Spec.Networking.Pods)
		assert.Equal(t, runtimeShoot.Networking.Services, *shoot.Spec.Networking.Services)
	})
}
