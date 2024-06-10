package extender

import (
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegionExtender(t *testing.T) {
	t.Run("Use region from Runtime", func(t *testing.T) {
		// given
		runtime := imv1.Runtime{
			Spec: imv1.RuntimeSpec{
				Shoot: imv1.RuntimeShoot{
					Region:         "ap-northeast-1",
					PlatformRegion: "cf-asia",
				},
			},
		}
		shoot := fixEmptyGardenerShoot("shoot", "kcp-system")

		// when
		err := ExtendWithRegion(runtime, &shoot)
		require.NoError(t, err)

		// then
		assert.Equal(t, "ap-northeast-1", shoot.Spec.Region)
	})

	t.Run("Use default region for EU access", func(t *testing.T) {
		// given
		runtime := imv1.Runtime{
			Spec: imv1.RuntimeSpec{
				Shoot: imv1.RuntimeShoot{
					Region:         "ap-northeast-1",
					PlatformRegion: "cf-eu11",
				},
			},
		}
		shoot := fixEmptyGardenerShoot("shoot", "kcp-system")

		// when
		err := ExtendWithRegion(runtime, &shoot)
		require.NoError(t, err)

		// then
		assert.Equal(t, "eu-central-1", shoot.Spec.Region)
	})
}
