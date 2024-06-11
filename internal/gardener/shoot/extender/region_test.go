package extender

import (
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegionExtender(t *testing.T) {
	for _, testCase := range []struct {
		name           string
		region         string
		platformRegion string
		providerType   string
		expectedRegion string
	}{
		{
			name:           "Use region from runtime for AWS",
			region:         "ap-northeast-1",
			platformRegion: "cf-asia",
			providerType:   "aws",
			expectedRegion: "ap-northeast-1",
		},
		{
			name:           "Use region from runtime for Azure",
			region:         "southeastasia",
			platformRegion: "cf-asia",
			providerType:   "aws",
			expectedRegion: "southeastasia",
		},
		{
			name:           "Use region from runtime for GCP",
			region:         "asia-northeast2",
			platformRegion: "cf-asia",
			providerType:   "gcp",
			expectedRegion: "asia-northeast2",
		},
		{
			name:           "Use region from runtime for Openstack",
			region:         "eu-de-1",
			platformRegion: "cf-eu11",
			providerType:   "openstack",
			expectedRegion: "eu-de-1",
		},
		{
			name:           "Replace region for EU Access on AWS",
			region:         "ap-northeast-1",
			platformRegion: "cf-eu11",
			providerType:   "aws",
			expectedRegion: "eu-central-1",
		},
		{
			name:           "Replace region for EU Access on Azure",
			region:         "uksouth",
			platformRegion: "cf-ch20",
			providerType:   "azure",
			expectedRegion: "switzerlandnorth",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			runtime := imv1.Runtime{
				Spec: imv1.RuntimeSpec{
					Shoot: imv1.RuntimeShoot{
						Region:         testCase.region,
						PlatformRegion: testCase.platformRegion,
						Provider: imv1.Provider{
							Type: testCase.providerType,
						},
					},
				},
			}
			shoot := fixEmptyGardenerShoot("shoot", "kcp-system")

			// when
			err := ExtendWithRegion(runtime, &shoot)
			require.NoError(t, err)

			// then
			assert.Equal(t, testCase.expectedRegion, shoot.Spec.Region)
		})
	}
}
