package shoot

import (
	"fmt"
	"strings"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/runtime"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/utilz"
	"github.com/onsi/gomega/types"
)

func NewProviderMatcher(v any, path string) types.GomegaMatcher {
	return &ProviderMatcher{
		toMatch:  v,
		rootPath: path,
	}
}

type ProviderMatcher struct {
	toMatch  interface{}
	fails    []string
	rootPath string
}

func (m *ProviderMatcher) getPath(p string) string {
	return fmt.Sprintf("%s/%s", m.rootPath, p)
}

func (m *ProviderMatcher) Match(actual interface{}) (success bool, err error) {
	aProvider, err := utilz.Get[v1beta1.Provider](actual)
	if err != nil {
		return false, err
	}

	eProvider, err := utilz.Get[v1beta1.Provider](m.toMatch)
	if err != nil {
		return false, err
	}

	for _, matcher := range []propertyMatcher{
		{
			path:          m.getPath("controlPlaneConfig"),
			GomegaMatcher: runtime.NewRawExtensionMatcher(eProvider.ControlPlaneConfig),
			expected:      aProvider.ControlPlaneConfig,
		},
		{
			path:          m.getPath("infrastructureConfig"),
			GomegaMatcher: runtime.NewRawExtensionMatcher(eProvider.InfrastructureConfig),
			expected:      aProvider.InfrastructureConfig,
		},
	} {
		ok, err := matcher.Match(matcher.expected)
		if err != nil {
			return false, err
		}

		if !ok {
			msg := matcher.FailureMessage(matcher.expected)
			if matcher.path != "" {
				msg = fmt.Sprintf("%s: %s", matcher.path, msg)
			}
			m.fails = append(m.fails, msg)
		}
	}

	return len(m.fails) == 0, nil
}

func (m *ProviderMatcher) NegatedFailureMessage(_ interface{}) string {
	return "expected should not equal actual"
}

func (m *ProviderMatcher) FailureMessage(_ interface{}) string {
	return strings.Join(m.fails, "\n")
}

type propertyMatcher = struct {
	types.GomegaMatcher
	path     string
	expected interface{}
}
