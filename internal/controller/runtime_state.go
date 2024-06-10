package controller

import (
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

// the state of controlled system (k8s cluster)
type systemState struct {
	instance imv1.Runtime
	snapshot imv1.RuntimeStatus
}

func (s *systemState) saveRuntimeStatus() {
	result := s.instance.Status.DeepCopy()
	if result == nil {
		result = &imv1.RuntimeStatus{}
	}
	s.snapshot = *result
}
