package fsm

import (
	"context"
	"fmt"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnWaitForShootReconcile(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Waiting for shoot reconcile state")

	lastOperation := s.shoot.Status.LastOperation

	switch s.shoot.Status.LastOperation.State {
	case gardener.LastOperationStateProcessing, gardener.LastOperationStatePending:
		msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", s.shoot.Name, lastOperation.State)
		m.log.Info(msg)

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonProcessing,
			"Unknown",
			"Shoot update is in progress")

		return updateStatusAndRequeueAfter(gardenerRequeueDuration)

	case gardener.LastOperationStateSucceeded:
		if !s.instance.IsStateWithConditionAndStatusSet(imv1.RuntimeStatePending, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonProcessing, "True") {
			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonProcessing,
				"True",
				"Shoot update is completed")

			return updateStatusAndRequeue()
		}

		msg := fmt.Sprintf("Shoot %s successfully updated, moving to processing", s.shoot.Name)
		m.log.Info(msg)

		return switchState(sFnProcessShoot)

	case gardener.LastOperationStateFailed:
		var reason ErrReason

		if len(s.shoot.Status.LastErrors) > 0 {
			reason = gardenerErrCodesToErrReason(s.shoot.Status.LastErrors...)
		}

		msg := fmt.Sprintf("error during cluster processing: reconcilation failed for shoot %s, reason: %s, exiting with no retry", s.shoot.Name, reason)
		m.log.Info(msg)

		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonProcessingErr,
			"False",
			string(reason))

		return updateStatusAndStop()
	}

	m.log.Info("Update did not processed, exiting with no retry")
	return stop()
}
