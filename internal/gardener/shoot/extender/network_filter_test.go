package extender

import (
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkingFilterExtender(t *testing.T) {
	t.Run("Create networking-filter extension", func(t *testing.T) {
		// given
		runtimeShoot := imv1.Runtime{
			Spec: imv1.RuntimeSpec{
				Shoot: imv1.RuntimeShoot{
					Name: "myshoot",
				},
			},
		}
		shoot := fixEmptyGardenerShoot("test", "dev")

		// when
		err := ExtendWithNetworkFilter(runtimeShoot, &shoot)

		// then
		require.NoError(t, err)
		assert.Equal(t, false, *shoot.Spec.Extensions[0].Disabled)
		assert.Equal(t, NetworkFilterType, shoot.Spec.Extensions[0].Type)
	})
}
