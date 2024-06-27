package fsm

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/types"
)

func haveName(name string) types.GomegaMatcher {
	return &stateFnNameMatcher{
		expName: name,
	}
}

type stateFnNameMatcher struct {
	actName string
	expName string
}

func lastIndexOf(s string, of rune) int {
	result := -1
	for i, r := range s {
		if r == of {
			result = i
		}
	}
	return result
}

func (m *stateFnNameMatcher) Match(actual any) (success bool, err error) {
	actualFn, ok := actual.(stateFn)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects stateFn")
	}

	r := regexp.MustCompile(".func[0-9]*|.glob")

	m.actName = actualFn.name()
	m.actName = r.ReplaceAllString(m.actName, "")

	aliof := lastIndexOf(m.actName, '.')
	if aliof != -1 {
		m.actName = m.actName[aliof+1:]
	}

	return m.actName == m.expName, nil
}

func (m *stateFnNameMatcher) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nto be equal to\n\t%s", m.actName, m.expName)
}

func (m *stateFnNameMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nnot to be equal to\n\t%s", m.actName, m.expName)
}
