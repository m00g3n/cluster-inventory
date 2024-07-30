package extender

import (
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTolerationsExtender(t *testing.T) {
	t.Run("Extend Tolerations", func(t *testing.T) {
		// given
		shoot := fixEmptyGardenerShoot("shoot", "kcp-system")
		basicRuntime := imv1.Runtime{
			ObjectMeta: v1.ObjectMeta{
				Name:      "runtime",
				Namespace: "namespace",
			},
			Spec: imv1.RuntimeSpec{
				Shoot: imv1.RuntimeShoot{
					Region: "me-central2",
				},
			},
		}

		// when
		err := ExtendWithTolerations(basicRuntime, &shoot)
		require.NoError(t, err)

		// then
		expectedTolerations := []gardener.Toleration{
			{
				Key: "ksa-assured-workload",
			},
		}
		assert.Equal(t, expectedTolerations, shoot.Spec.Tolerations)
	})
}
