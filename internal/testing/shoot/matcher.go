package shoot

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"sigs.k8s.io/yaml"
)

var (
	errInvalidType = fmt.Errorf("invalid type")
)

type Matcher struct {
	toMatch interface{}
	fails   []string
}

func NewMatcher(i interface{}) types.GomegaMatcher {
	return &Matcher{
		toMatch: i,
	}
}

func getShoot(i interface{}) (shoot v1beta1.Shoot, err error) {
	if i == nil {
		return v1beta1.Shoot{}, fmt.Errorf("invalid value nil")
	}

	switch v := i.(type) {
	case string:
		err = yaml.Unmarshal([]byte(v), &shoot)
		return shoot, err

	case v1beta1.Shoot:
		return v, nil

	case *v1beta1.Shoot:
		return *v, nil

	default:
		return v1beta1.Shoot{}, fmt.Errorf(`%w: %s`, errInvalidType, reflect.TypeOf(v))
	}
}

type matcher struct {
	types.GomegaMatcher
	path   string
	actual interface{}
}

func (m *Matcher) Match(actual interface{}) (success bool, err error) {
	aShoot, err := getShoot(actual)
	if err != nil {
		return false, err
	}

	eShoot, err := getShoot(m.toMatch)
	if err != nil {
		return false, err
	}

	for _, matcher := range []matcher{
		{
			GomegaMatcher: gomega.Equal(eShoot.TypeMeta),
			actual:        aShoot.TypeMeta,
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.Name),
			actual:        aShoot.Name,
			path:          "metadata/name",
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.Namespace),
			actual:        aShoot.Namespace,
			path:          "metadata/namespace",
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.Labels),
			actual:        aShoot.Labels,
			path:          "metadata/labels",
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.Annotations),
			actual:        aShoot.Annotations,
			path:          "metadata/annotations",
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.OwnerReferences),
			actual:        aShoot.OwnerReferences,
			path:          "metadata/ownerReferences",
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.Finalizers),
			actual:        aShoot.Finalizers,
			path:          "metadata/finalizers",
		},
		{
			GomegaMatcher: gomega.Equal(eShoot.Spec),
			actual:        aShoot.Spec,
			path:          "spec",
		},
	} {
		ok, err := matcher.Match(matcher.actual)
		if err != nil {
			return false, err
		}

		if !ok {
			msg := matcher.FailureMessage(matcher.actual)
			if matcher.path != "" {
				msg = fmt.Sprintf("%s: %s", matcher.path, msg)
			}
			m.fails = append(m.fails, msg)
		}
	}

	return len(m.fails) == 0, nil
}

func (m *Matcher) NegatedFailureMessage(_ interface{}) string {
	return "expected should not equal actual"
}

func (m *Matcher) FailureMessage(_ interface{}) string {
	return strings.Join(m.fails, "\n")
}
