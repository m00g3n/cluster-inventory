package shoot

import (
	"crypto/md5"
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
	path     string
	expected interface{}
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

	hashedExpectedControlPlaneConfig := md5.Sum(eShoot.Spec.Provider.ControlPlaneConfig.Raw)
	hashedActualControlPlaneConfig := md5.Sum(aShoot.Spec.Provider.ControlPlaneConfig.Raw)

	hashedExpectedInfrastructureConfig := md5.Sum(eShoot.Spec.Provider.InfrastructureConfig.Raw)
	hashedActualInfrastructureConfig := md5.Sum(aShoot.Spec.Provider.InfrastructureConfig.Raw)

	// Note: we define separate matchers for each field to make input more readable
	// Annotations are not matched as they are not relevant for the comparison ; both KIM, and Provisioner have different set of annotations
	for _, matcher := range []matcher{
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
			GomegaMatcher: gomega.HaveKey("account"),
			expected:      aShoot.Labels,
			path:          "metadata/labels/account",
		},
		{
			GomegaMatcher: gomega.HaveKey("subaccount"),
			expected:      aShoot.Labels,
			path:          "metadata/labels/subaccount",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.CloudProfileName),
			expected:      aShoot.Spec.CloudProfileName,
			path:          "spec/CloudProfileName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.DNS),
			expected:      aShoot.Spec.DNS,
			path:          "spec/DNS",
		},
		{
			GomegaMatcher: NewExtensionMatcher(eShoot.Spec.Extensions),
			expected:      aShoot.Spec.Extensions,
			path:          "spec/Extensions",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Hibernation),
			expected:      aShoot.Spec.Hibernation,
			path:          "spec/Hibernation",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Kubernetes.Version),
			expected:      aShoot.Spec.Kubernetes.Version,
			path:          "spec/Kubernetes.Version",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Kubernetes.EnableStaticTokenKubeconfig),
			expected:      aShoot.Spec.Kubernetes.EnableStaticTokenKubeconfig,
			path:          "spec/Kubernetes.EnableStaticTokenKubeconfig",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Networking.Type),
			expected:      aShoot.Spec.Networking.Type,
			path:          "spec/Networking.Type",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Networking.Nodes),
			expected:      aShoot.Spec.Networking.Nodes,
			path:          "spec/Networking.Nodes",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Networking.Pods),
			expected:      aShoot.Spec.Networking.Pods,
			path:          "spec/Networking.Pods",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Networking.Services),
			expected:      aShoot.Spec.Networking.Services,
			path:          "spec/Networking.Services",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Maintenance.AutoUpdate.KubernetesVersion),
			expected:      aShoot.Spec.Maintenance.AutoUpdate.KubernetesVersion,
			path:          "spec/Maintenance.AutoUpdate.KubernetesVersion",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Maintenance.AutoUpdate.MachineImageVersion),
			expected:      aShoot.Spec.Maintenance.AutoUpdate.MachineImageVersion,
			path:          "spec/Maintenance.AutoUpdate.MachineImageVersion",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Monitoring),
			expected:      aShoot.Spec.Monitoring,
			path:          "spec/Monitoring",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Provider.Type),
			expected:      aShoot.Spec.Provider.Type,
			path:          "spec/Provider.Type",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Provider.Workers),
			expected:      aShoot.Spec.Provider.Workers,
			path:          "spec/Provider.Workers",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(hashedExpectedControlPlaneConfig),
			expected:      hashedActualControlPlaneConfig,
			path:          "spec/Provider.ControlPlaneConfig",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(hashedExpectedInfrastructureConfig),
			expected:      hashedActualInfrastructureConfig,
			path:          "spec/Provider.InfrastructureConfig",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Purpose),
			expected:      aShoot.Spec.Purpose,
			path:          "spec/Purpose",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Region),
			expected:      aShoot.Spec.Region,
			path:          "spec/Region",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.SecretBindingName),
			expected:      aShoot.Spec.SecretBindingName,
			path:          "spec/SecretBindingName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.Tolerations),
			expected:      aShoot.Spec.Tolerations,
			path:          "spec/Tolerations",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.ExposureClassName),
			expected:      aShoot.Spec.ExposureClassName,
			path:          "spec/ExposureClassName",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.ControlPlane),
			expected:      aShoot.Spec.ControlPlane,
			path:          "spec/ControlPlane",
		},
		{
			GomegaMatcher: gomega.BeComparableTo(eShoot.Spec.CredentialsBindingName),
			expected:      aShoot.Spec.CredentialsBindingName,
			path:          "spec/CredentialsBindingName",
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
