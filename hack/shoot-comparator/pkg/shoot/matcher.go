package shoot

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/errors"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
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

	matchers := []propertyMatcher{
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
			GomegaMatcher: gstruct.MatchElements(idExtension, gstruct.IgnoreExtras, extensions(aShoot.Spec.Extensions)),
			expected:      eShoot.Spec.Extensions,
			path:          "spec/extensions",
		},
		{
			GomegaMatcher: gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Version":                     gomega.BeComparableTo(aShoot.Spec.Kubernetes.Version),
				"EnableStaticTokenKubeconfig": gomega.BeComparableTo(aShoot.Spec.Kubernetes.EnableStaticTokenKubeconfig),
				"KubeAPIServer":               newKubeAPIServerMatcher(aShoot.Spec.Kubernetes),
			}),
			expected: eShoot.Spec.Kubernetes,
			path:     "spec/kubernetes",
		},
		{
			GomegaMatcher: newNetworkingMatcher(aShoot.Spec),
			expected:      eShoot.Spec.Networking,
			path:          "spec/networking",
		},
		{
			GomegaMatcher: newMaintenanceMatcher(eShoot.Spec),
			expected:      aShoot.Spec.Maintenance,
			path:          "spec/maintenance",
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
			GomegaMatcher: newDNSMatcher(aShoot.Spec.DNS),
			path:          "spec/dns",
			expected:      eShoot.Spec.DNS,
		},
		{
			GomegaMatcher: gstruct.MatchElements(
				idToleration,
				gstruct.IgnoreExtras,
				tolerations(aShoot.Spec.Tolerations),
			),
			expected: eShoot.Spec.Tolerations,
			path:     "spec/tolerations",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.ExposureClassName),
			expected:      aShoot.Spec.ExposureClassName,
			path:          "spec/exposureClassName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.ControlPlane),
			expected:      aShoot.Spec.ControlPlane,
			path:          "spec/controlPlane",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.CloudProfile),
			expected:      aShoot.Spec.CloudProfile,
			path:          "spec/cloudProfile",
		},
		{
			GomegaMatcher: NewProviderMatcher(eShoot.Spec.Provider, "spec/provider"),
			expected:      aShoot.Spec.Provider,
			path:          "spec/provider",
		},
		{
			GomegaMatcher: gomega.SatisfyAll(mapMatchers(aShoot.Labels)...),
			expected:      eShoot.Labels,
			path:          "metadata/labels",
		},
		{
			GomegaMatcher: gomega.SatisfyAll(mapMatchers(aShoot.Annotations)...),
			expected:      eShoot.Annotations,
			path:          "metadata/annotations",
		},
	}

	for _, matcher := range matchers {
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

func mapMatchers(m map[string]string) []types.GomegaMatcher {
	mLen := len(m)
	if mLen == 0 {
		return []types.GomegaMatcher{
			gomega.BeEmpty(),
		}
	}

	out := make([]types.GomegaMatcher, mLen)
	index := 0
	for key, val := range m {
		matcher := gomega.HaveKeyWithValue(key, val)
		out[index] = matcher
		index++
	}
	return out
}

func val(v interface{}) string {
	if reflect.ValueOf(v).Kind() == reflect.Pointer && reflect.ValueOf(v).IsNil() {
		return ""
	}

	if reflect.ValueOf(v).Kind() == reflect.Pointer {
		return fmt.Sprintf("%v", reflect.ValueOf(v).Elem())
	}

	return fmt.Sprintf("%v", v)
}

func idToleration(v interface{}) string {
	toleration, ok := v.(v1beta1.Toleration)
	if !ok {
		panic("invalid type")
	}
	return fmt.Sprintf("%s:%s", toleration.Key, val(toleration.Value))
}

func tolerations(ts []v1beta1.Toleration) gstruct.Elements {
	out := map[string]types.GomegaMatcher{}
	for _, t := range ts {
		ID := idToleration(t)
		out[ID] = gstruct.MatchAllFields(gstruct.Fields{
			"Key":   gomega.BeComparableTo(t.Key),
			"Value": gomega.BeComparableTo(t.Value),
		})
	}
	return out
}

func idProvider(v interface{}) string {
	provider, ok := v.(v1beta1.DNSProvider)
	if !ok {
		panic("invalid type")
	}

	return fmt.Sprintf("%s:%s:%s",
		val(provider.Type),
		val(provider.SecretName),
		val(provider.Primary))
}

func providers(ps []v1beta1.DNSProvider) gstruct.Elements {
	out := map[string]types.GomegaMatcher{}
	for _, p := range ps {
		ID := idProvider(p)

		domainsMatcher := gomega.BeNil()
		if p.Domains != nil {
			domainsMatcher = gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Include": gomega.ContainElements(p.Domains.Include),
			}))
		}

		out[ID] = gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Primary":    gomega.Equal(p.Primary),
			"SecretName": gomega.Equal(p.SecretName),
			"Type":       gomega.Equal(p.Type),
			"Domains":    domainsMatcher,
		})
	}

	return out
}

func newDNSMatcher(dns *v1beta1.DNS) types.GomegaMatcher {
	if dns == nil {
		return gomega.BeNil()
	}

	return gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"Domain": gomega.BeComparableTo(dns.Domain),
		"Providers": gstruct.MatchElements(idProvider, gstruct.IgnoreExtras,
			providers(dns.Providers)),
	}))
}

func newMaintenanceMatcher(spec v1beta1.ShootSpec) types.GomegaMatcher {
	if spec.Maintenance == nil {
		return gomega.BeNil()
	}

	return gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"AutoUpdate": gomega.BeComparableTo(spec.Maintenance.AutoUpdate),
	}))
}

func newNetworkingMatcher(spec v1beta1.ShootSpec) types.GomegaMatcher {
	if spec.Networking == nil {
		return gomega.BeNil()
	}

	return gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
		"Type":     gomega.BeComparableTo(spec.Networking.Type),
		"Nodes":    gomega.BeComparableTo(spec.Networking.Nodes),
		"Pods":     gomega.BeComparableTo(spec.Networking.Pods),
		"Services": gomega.BeComparableTo(spec.Networking.Services),
	}))
}

func newKubeAPIServerMatcher(k v1beta1.Kubernetes) types.GomegaMatcher {
	if k.KubeAPIServer == nil {
		return gomega.BeNil()
	}

	return gstruct.PointTo(gstruct.MatchFields(
		gstruct.IgnoreExtras,
		gstruct.Fields{
			"OIDCConfig": newOIDCConfigMatcher(k.KubeAPIServer),
		},
	))
}

func newOIDCConfigMatcher(c *v1beta1.KubeAPIServerConfig) types.GomegaMatcher {
	if c == nil || c.OIDCConfig == nil {
		return gomega.BeNil()
	}

	return gstruct.PointTo(gstruct.MatchFields(
		gstruct.IgnoreExtras,
		gstruct.Fields{
			"CABundle":       gomega.BeComparableTo(c.OIDCConfig.CABundle),
			"ClientID":       gomega.BeComparableTo(c.OIDCConfig.ClientID),
			"GroupsClaim":    gomega.BeComparableTo(c.OIDCConfig.GroupsClaim),
			"GroupsPrefix":   gomega.BeComparableTo(c.OIDCConfig.GroupsPrefix),
			"IssuerURL":      gomega.BeComparableTo(c.OIDCConfig.IssuerURL),
			"RequiredClaims": gomega.BeComparableTo(c.OIDCConfig.RequiredClaims),
			"SigningAlgs":    gomega.ContainElements(c.OIDCConfig.SigningAlgs),
			"UsernameClaim":  gomega.BeComparableTo(c.OIDCConfig.UsernameClaim),
			"UsernamePrefix": gomega.BeComparableTo(c.OIDCConfig.UsernamePrefix),
		},
	))
}

func idExtension(v interface{}) string {
	e, ok := v.(v1beta1.Extension)
	if !ok {
		panic("invalid type")
	}

	return e.Type
}

func extensions(es []v1beta1.Extension) gstruct.Elements {
	out := map[string]types.GomegaMatcher{}
	for _, e := range es {
		ID := idExtension(e)
		out[ID] = gstruct.MatchAllFields(gstruct.Fields{
			"Type":           gomega.BeComparableTo(e.Type),
			"ProviderConfig": gomega.BeComparableTo(e.ProviderConfig),
			"Disabled":       gomega.BeComparableTo(e.Disabled),
		})
	}
	return out
}
