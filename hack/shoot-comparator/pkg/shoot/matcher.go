package shoot

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/errors"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"sigs.k8s.io/yaml"
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
		return v1beta1.Shoot{}, fmt.Errorf(`%w: %s`, errors.ErrInvalidType, reflect.TypeOf(v))
	}
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

	// Note: we define separate matchers for each field to make input more readable
	// Annotations are not matched as they are not relevant for the comparison ; both KIM, and Provisioner have different set of annotations
	for _, matcher := range []propertyMatcher{
		// We need to skip comparing type meta as Provisioner doesn't set it.
		// It is simpler to skip it than to make fix in the Provisioner.
		//{
		//	GomegaMatcher: gomega.BeComparableTo(eShoot.TypeMeta),
		//	expected:      aShoot.TypeMeta,
		//},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Name),
			expected:      aShoot.Name,
			path:          "metadata/name",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Namespace),
			expected:      aShoot.Namespace,
			path:          "metadata/namespace",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Labels),
			expected:      aShoot.Labels,
			path:          "metadata/labels",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Addons),
			expected:      aShoot.Spec.Addons,
			path:          "spec/addons",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.CloudProfileName),
			expected:      aShoot.Spec.CloudProfileName,
			path:          "spec/cloudProfileName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.DNS),
			expected:      aShoot.Spec.DNS,
			path:          "spec/dns",
		},
		{
			GomegaMatcher: NewExtensionMatcher(eShoot.Spec.Extensions),
			expected:      aShoot.Spec.Extensions,
			path:          "spec/extensions",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Hibernation),
			expected:      aShoot.Spec.Hibernation,
			path:          "spec/hibernation",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Kubernetes),
			expected:      aShoot.Spec.Kubernetes,
			path:          "spec/kubernetes",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Networking),
			expected:      aShoot.Spec.Networking,
			path:          "spec/networking",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Maintenance),
			expected:      aShoot.Spec.Maintenance,
			path:          "spec/maintenance",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Monitoring),
			expected:      aShoot.Spec.Monitoring,
			path:          "spec/monitoring",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Purpose),
			expected:      aShoot.Spec.Purpose,
			path:          "spec/purpose",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Region),
			expected:      aShoot.Spec.Region,
			path:          "spec/region",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.SecretBindingName),
			expected:      aShoot.Spec.SecretBindingName,
			path:          "spec/secretBindingName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.SeedName),
			expected:      aShoot.Spec.SeedName,
			path:          "spec/seedName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.SeedSelector),
			expected:      aShoot.Spec.SeedSelector,
			path:          "spec/seedSelector",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Resources),
			expected:      aShoot.Spec.Resources,
			path:          "spec/resources",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Tolerations),
			expected:      aShoot.Spec.Tolerations,
			path:          "spec/tolerations",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.ExposureClassName),
			expected:      aShoot.Spec.ExposureClassName,
			path:          "spec/exposureClassName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.SystemComponents),
			expected:      aShoot.Spec.SystemComponents,
			path:          "spec/systemComponents",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.ControlPlane),
			expected:      aShoot.Spec.ControlPlane,
			path:          "spec/controlPlane",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.SchedulerName),
			expected:      aShoot.Spec.SchedulerName,
			path:          "spec/schedulerName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.CloudProfile),
			expected:      aShoot.Spec.CloudProfile,
			path:          "spec/cloudProfile",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.CredentialsBindingName),
			expected:      aShoot.Spec.CredentialsBindingName,
			path:          "spec/credentialsBindingName",
		},
		{
			GomegaMatcher: NewProviderMatcher(eShoot.Spec.Provider, "spec/provider"),
			expected:      aShoot.Spec.Provider,
			path:          "spec/provider",
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

func (m *Matcher) NegatedFailureMessage(_ interface{}) string {
	return "expected should not equal actual"
}

func (m *Matcher) FailureMessage(_ interface{}) string {
	return strings.Join(m.fails, "\n")
}
