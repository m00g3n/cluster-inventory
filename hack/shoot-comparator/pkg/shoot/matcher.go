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
	shootActual, err := getShoot(actual)
	if err != nil {
		return false, err
	}

	shootToMatch, err := getShoot(m.toMatch)
	if err != nil {
		return false, err
	}

	matchers := []propertyMatcher{
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Name),
			actual:        shootActual.Name,
			path:          "metadata/name",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Namespace),
			actual:        shootActual.Namespace,
			path:          "metadata/namespace",
		},
		{
			GomegaMatcher: gstruct.MatchElements(idExtension, gstruct.IgnoreExtras, extensions(shootActual.Spec.Extensions)),
			actual:        shootToMatch.Spec.Extensions,
			path:          "spec/extensions",
		},
		{
			GomegaMatcher: gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Version":                     gomega.BeComparableTo(shootActual.Spec.Kubernetes.Version),
				"EnableStaticTokenKubeconfig": gomega.BeComparableTo(shootActual.Spec.Kubernetes.EnableStaticTokenKubeconfig),
				"KubeAPIServer":               newKubeAPIServerMatcher(shootActual.Spec.Kubernetes),
			}),
			actual: shootToMatch.Spec.Kubernetes,
			path:   "spec/kubernetes",
		},
		{
			GomegaMatcher: newNetworkingMatcher(shootActual.Spec),
			actual:        shootToMatch.Spec.Networking,
			path:          "spec/networking",
		},
		{
			GomegaMatcher: newMaintenanceMatcher(shootToMatch.Spec),
			actual:        shootActual.Spec.Maintenance,
			path:          "spec/maintenance",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Spec.Purpose),
			actual:        shootActual.Spec.Purpose,
			path:          "spec/purpose",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Spec.Region),
			actual:        shootActual.Spec.Region,
			path:          "spec/region",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Spec.SecretBindingName),
			actual:        shootActual.Spec.SecretBindingName,
			path:          "spec/secretBindingName",
		},
		{
			GomegaMatcher: newDNSMatcher(shootActual.Spec.DNS),
			path:          "spec/dns",
			actual:        shootToMatch.Spec.DNS,
		},
		{
			GomegaMatcher: gstruct.MatchElements(
				idToleration,
				gstruct.IgnoreExtras,
				tolerations(shootActual.Spec.Tolerations),
			),
			actual: shootToMatch.Spec.Tolerations,
			path:   "spec/tolerations",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Spec.ExposureClassName),
			actual:        shootActual.Spec.ExposureClassName,
			path:          "spec/exposureClassName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Spec.ControlPlane),
			actual:        shootActual.Spec.ControlPlane,
			path:          "spec/controlPlane",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(shootToMatch.Spec.CloudProfile),
			actual:        shootActual.Spec.CloudProfile,
			path:          "spec/cloudProfile",
		},
		{
			GomegaMatcher: NewProviderMatcher(shootToMatch.Spec.Provider, "spec/provider"),
			actual:        shootActual.Spec.Provider,
			path:          "spec/provider",
		},
		{
			GomegaMatcher: gomega.SatisfyAll(mapMatchers(shootActual.Labels)...),
			actual:        shootToMatch.Labels,
			path:          "metadata/labels",
		},
		{
			GomegaMatcher: gomega.SatisfyAll(mapMatchers(shootActual.Annotations)...),
			actual:        shootToMatch.Annotations,
			path:          "metadata/annotations",
		},
	}

	for _, matcher := range matchers {
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
