package fsm

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// to save the runtime status at the begining of the reconciliation
func sFnTakeSnapshot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	s.saveRuntimeStatus()
	s.shoot = nil

	shoot, err := m.ShootClient.Get(ctx, s.instance.Name, v1.GetOptions{})

	if err != nil {
		if !apierrors.IsNotFound(err) {
			m.log.Info("Failed to get Gardener shoot", "error", err)
			return updateStatusAndRequeue()
		}
	} else if shoot != nil {
		s.shoot = shoot.DeepCopy()
	}

	return switchState(sFnInitialize)
}
