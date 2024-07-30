package v1_test

import (
	"fmt"
	"testing"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Runtime_IsControlledByProvisioner(t *testing.T) {
	var tests = []struct {
		desc     string
		rt       imv1.Runtime
		expected bool
	}{
		{
			desc:     "rt without labels (zero value) should be controlled by the provisioner",
			rt:       imv1.Runtime{},
			expected: true,
		},
		{
			desc: "rt without labels (empty map) should be controlled by the provisioner",
			rt: imv1.Runtime{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{},
				},
			},
			expected: true,
		},
		{
			desc: fmt.Sprintf(`rt without label: "%s" should be controlled by the provisioner`, imv1.LabelControlledByProvisioner),
			rt: imv1.Runtime{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"test": "me"},
				},
			},
			expected: true,
		},
		{
			desc: fmt.Sprintf(`rt with label: "%s=true" should be controlled by the provisioner`, imv1.LabelControlledByProvisioner),
			rt: imv1.Runtime{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{imv1.LabelControlledByProvisioner: "true"},
				},
			},
			expected: true,
		},
		{
			desc: fmt.Sprintf(`rt with label: "%s=sth" should be controlled by the provisioner`, imv1.LabelControlledByProvisioner),
			rt: imv1.Runtime{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{imv1.LabelControlledByProvisioner: "sth"},
				},
			},
			expected: true,
		},
		{
			desc: fmt.Sprintf(`rt with label: "%s=false" should NOT be controlled by the provisioner`, imv1.LabelControlledByProvisioner),
			rt: imv1.Runtime{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{imv1.LabelControlledByProvisioner: "false"},
				},
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			actual := tt.rt.IsControlledByProvisioner()
			assert.Equal(t, tt.expected, actual, tt.desc)
		})
	}
}
