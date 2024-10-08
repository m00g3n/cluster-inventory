package runtime

import (
	"sort"

	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/utilz"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/runtime"
)

type RawExtensionMatcher struct {
	toMatch interface{}
}

func NewRawExtensionMatcher(v any) types.GomegaMatcher {
	return &RawExtensionMatcher{
		toMatch: v,
	}
}

func (m *RawExtensionMatcher) Match(actual interface{}) (bool, error) {
	if actual == nil && m.toMatch == nil {
		return true, nil
	}

	aRawExtension, err := utilz.Get[runtime.RawExtension](actual)
	if err != nil {
		return false, err
	}

	eRawExtension, err := utilz.Get[runtime.RawExtension](m.toMatch)
	if err != nil {
		return false, err
	}

	sort.Sort(sortBytes(aRawExtension.Raw))
	sort.Sort(sortBytes(eRawExtension.Raw))

	return gomega.BeComparableTo(aRawExtension.Raw).Match(eRawExtension.Raw)
}

func (m *RawExtensionMatcher) NegatedFailureMessage(_ interface{}) string {
	return "expected should not equal actual"
}

func (m *RawExtensionMatcher) FailureMessage(_ interface{}) string {
	return "expected should equal actual"
}

type sortBytes []byte

func (s sortBytes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortBytes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortBytes) Len() int {
	return len(s)
}
