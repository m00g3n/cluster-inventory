package extender

import (
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnnotationsAndLabelsExtender(t *testing.T) {
	licenceType := "licence"

	for _, testCase := range []struct {
		name                string
		runtime             imv1.Runtime
		expectedAnnotations map[string]string
		expectedLabels      map[string]string
	}{
		{
			name: "Extend with basic annotations",
			runtime: imv1.Runtime{
				ObjectMeta: v1.ObjectMeta{
					Name:      "runtime",
					Namespace: "namespace",
					Labels: map[string]string{
						"kyma-project.io/runtime-id": "runtime-id",
					},
				},
			},
			expectedAnnotations: map[string]string{
				"infrastructuremanager.kyma-project.io/runtime-id": "runtime-id"},
			expectedLabels: map[string]string{
				"account":    "account-id",
				"subaccount": "subaccount-id"},
		},
		{
			name: "Extend with licence type annotation",
			runtime: imv1.Runtime{
				ObjectMeta: v1.ObjectMeta{
					Name:      "runtime",
					Namespace: "namespace",
					Labels: map[string]string{
						"kyma-project.io/runtime-id": "runtime-id",
					},
				},
				Spec: imv1.RuntimeSpec{
					Shoot: imv1.RuntimeShoot{
						LicenceType: &licenceType,
					},
				},
			},
			expectedAnnotations: map[string]string{
				"infrastructuremanager.kyma-project.io/runtime-id":   "runtime-id",
				"infrastructuremanager.kyma-project.io/licence-type": "licence"},
			expectedLabels: map[string]string{
				"account":    "account-id",
				"subaccount": "subaccount-id"},
		},
	} {
		// given
		shoot := fixEmptyGardenerShoot("shoot", "kcp-system")

		// when
		err := ExtendWithAnnotationsAndLabels(testCase.runtime, &shoot)
		require.NoError(t, err)

		// then
		assert.Equal(t, testCase.expectedAnnotations, shoot.Annotations)
	}
}
