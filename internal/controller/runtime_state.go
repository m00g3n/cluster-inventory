package controller

import (
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// the state of controlled system (k8s cluster)
type systemState struct {
	instance imv1.Runtime

	objs []unstructured.Unstructured

	snapshot imv1.RuntimeStatus

	domainName string
}

func (s *systemState) saveRuntimeStatus() {
	result := s.instance.Status.DeepCopy()
	if result == nil {
		result = &imv1.RuntimeStatus{}
	}
	s.snapshot = *result
}
