package runtime

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

func loadTestdata(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	require.Greater(t, len(data), 0, "test data is empty")
	return data
}

func TestMatcher(t *testing.T) {
	testCases := []struct {
		actual   interface{}
		expected interface{}
		isOK     bool
		hasErr   bool
	}{
		{
			actual:   string(loadTestdata(t, "testdata/infra_cfg_azure_11.yaml")),
			expected: string(loadTestdata(t, "testdata/infra_cfg_azure_12.yaml")),
			isOK:     true,
			hasErr:   false,
		},
		{
			actual:   runtime.RawExtension{},
			expected: runtime.RawExtension{},
			isOK:     true,
			hasErr:   false,
		},
		{
			actual:   runtime.RawExtension{},
			expected: nil,
			isOK:     false,
			hasErr:   true,
		},
		{
			actual:   string(loadTestdata(t, "testdata/infra_cfg_azure_11.yaml")),
			expected: string(loadTestdata(t, "testdata/infra_cfg_azure_21.yaml")),
			isOK:     false,
			hasErr:   false,
		},
		{
			actual:   []byte("invalid type"),
			expected: []byte("invalid type"),
			isOK:     false,
			hasErr:   true,
		},
	}

	for _, tc := range testCases {
		matcher := NewRawExtensionMatcher(tc.actual)
		ok, err := matcher.Match(tc.expected)
		// THEN
		assert.Equal(t, tc.hasErr, err != nil, err)
		assert.Equal(t, tc.isOK, ok, matcher.FailureMessage(nil))
	}

}
