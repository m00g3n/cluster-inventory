package fsm

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
)

// to save the runtime status at the begining of the reconciliation
func sFnTakeSnapshot(_ context.Context, _ *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	s.saveRuntimeStatus()
	return switchState(sFnInitialize)
}
