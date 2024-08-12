package shoot

import (
	"fmt"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"reflect"
	"strings"
)

type ExtensionMatcher struct {
	toMatch interface{}
	fails   []string
}

func NewExtensionMatcher(i interface{}) types.GomegaMatcher {
	return &ExtensionMatcher{
		toMatch: i,
	}
}

func (m *ExtensionMatcher) Match(actual interface{}) (success bool, err error) {

	aExtensions, err := getExtension(actual)
	if err != nil {
		return false, err
	}

	eExtensions, err := getExtension(m.toMatch)
	if err != nil {
		return false, err
	}

	if len(aExtensions) != len(eExtensions) {
		m.fails = append(m.fails, "Extensions count mismatch")
		return false, nil
	}

	if len(aExtensions) == 0 && len(eExtensions) == 0 {
		return true, nil
	}

	findExtension := func(name string, extensions []v1beta1.Extension) v1beta1.Extension {
		for _, e := range extensions {
			if e.Type == name {
				return e
			}
		}

		return v1beta1.Extension{}
	}

	differenceFound := false

	for _, e := range eExtensions {
		a := findExtension(e.Type, aExtensions)
		if a.Type == "" {
			m.fails = append(m.fails, fmt.Sprintf("Extension %s not found in both expected and actual", e.Type))
			return false, nil
		}

		matcher := gomega.BeComparableTo(e)

		ok, err := matcher.Match(a)
		if err != nil {
			return false, err
		}

		if !ok {
			differenceFound = true
			m.fails = append(m.fails, matcher.FailureMessage(a))
		}
	}

	return !differenceFound, nil
}

func (m *ExtensionMatcher) NegatedFailureMessage(_ interface{}) string {
	return "expected should not equal actual"
}

func (m *ExtensionMatcher) FailureMessage(_ interface{}) string {
	return strings.Join(m.fails, "\n")
}

func getExtension(i interface{}) ([]v1beta1.Extension, error) {
	switch v := i.(type) {
	case []v1beta1.Extension:
		return v, nil

	default:
		return []v1beta1.Extension{}, fmt.Errorf(`%w: %s`, errInvalidType, reflect.TypeOf(v))
	}
}
