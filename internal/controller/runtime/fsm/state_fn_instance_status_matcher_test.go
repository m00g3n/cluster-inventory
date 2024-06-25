package fsm

/*
import (
	"fmt"

	. "github.com/onsi/gomega" //nolint:revive

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


type instanceConditionMatcher struct {
	actual   *metav1.Condition
	expected *metav1.Condition
}

func haveCondition(v *metav1.Condition) types.GomegaMatcher {
	return &instanceConditionMatcher{
		expected: v,
	}
}

func (m *instanceConditionMatcher) Match(actual any) (success bool, err error) {
	actualRt, ok := actual.(*imv1.Runtime)
	if !ok {
		return false, fmt.Errorf("instance condition matcher expects an *github.com/kyma-project/infrastructure-manager/api/v1.Runtime")
	}

	m.actual = meta.FindStatusCondition(actualRt.Status.Conditions, m.expected.Type)
	return Equal(m.expected).Match(m.actual)
}

func (m *instanceConditionMatcher) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto contain\n\t%s", m.actual, m.expected)
}

func (m *instanceConditionMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to contain\n\t%s", m.actual, m.expected)
}

*/
