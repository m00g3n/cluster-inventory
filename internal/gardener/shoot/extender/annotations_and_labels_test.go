package extender

import (
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// account: 77ddcd39-9a27-4f6e-9857-3680a976407c
// subaccount: cf725b0e-8784-4062-9e9f-b38aca762a87
// kcp.provisioner.kyma-project.io/runtime-id: b8a0d491-93d7-46d3-93e4-af692794a3be
// kcp.provisioner.kyma-project.io/licence-type

func TestAnnotationsAndLabelsExtender(t *testing.T) {
	for _, testCase := range []struct {
		name                string
		runtime             imv1.RuntimeShoot
		expectedAnnotations map[string]string
		expectedLabels      map[string]string
	}{
		{
			name: "Extend with basic annotations",
			expectedAnnotations: map[string]string{
				"infrastructuremanager.kyma-project.io/runtime-id": "runtime-id"},
			expectedLabels: map[string]string{
				"account":    "account-id",
				"subaccount": "subaccount-id"},
		},
		{
			name: "Extend with licence type annotation",
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
