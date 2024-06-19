package fsm

import (
	"context"
	"fmt"
	"strings"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardenerhelper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ErrReason string

func sFnPrepareCluster(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	if s.shoot == nil {
		panic("Fatal state machine logic problem: Shoot can never be nil in PrepareCluster state!")
	}

	if s.shoot.Spec.DNS == nil || s.shoot.Spec.DNS.Domain == nil {
		msg := fmt.Sprintf("DNS Domain is not set yet for shoot: %s, scheduling for retry", s.shoot.Name)
		m.log.Info(msg)
		return stopWithRequeue()
	}

	lastOperation := s.shoot.Status.LastOperation
	if lastOperation == nil {
		msg := fmt.Sprintf("Last operation is nil for shoot: %s, scheduling for retry", s.shoot.Name)
		m.log.Info(msg)
		return stopWithRequeue()
	}

	if lastOperation.Type == gardener.LastOperationTypeCreate {

		if lastOperation.State == gardener.LastOperationStateProcessing ||
			lastOperation.State == gardener.LastOperationStatePending {
			msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", s.shoot.Name, lastOperation.State)
			m.log.Info(msg)

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonShootCreationPending,
				"Unknown",
				"Shoot creation in progress")

			return stopWithRequeue()
		}

		if lastOperation.State == gardener.LastOperationStateSucceeded {
			msg := fmt.Sprintf("Shoot %s successfully created", s.shoot.Name)
			m.log.Info(msg)

			if !s.instance.IsRuntimeStateSet(imv1.RuntimeStatePending, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonShootCreationCompleted) {
				s.instance.UpdateStatePending(
					imv1.ConditionTypeRuntimeProvisioned,
					imv1.ConditionReasonShootCreationCompleted,
					"True",
					"Shoot creation completed")

				return stopWithRequeue()
			}

			return switchState(sFnProcessShoot)
		}

		if lastOperation.State == gardener.LastOperationStateFailed {
			if gardenerhelper.HasErrorCode(s.shoot.Status.LastErrors, gardener.ErrorInfraRateLimitsExceeded) {
				msg := fmt.Sprintf("Error during cluster provisioning: Rate limits exceeded for Shoot %s, scheduling for retry", s.shoot.Name)
				m.log.Info(msg)
				return stopWithRequeue()
			}

			msg := fmt.Sprintf("Provisioning failed for shoot: %s ! Last state: %s, Description: %s", s.shoot.Name, lastOperation.State, lastOperation.Description)
			m.log.Info(msg)

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonCreationError,
				"False",
				"Shoot creation failed")

			return stopWithNoRequeue()
		}
	}

	// Runtime update is in progress
	if lastOperation.Type == gardener.LastOperationTypeReconcile {

		if lastOperation.State == gardener.LastOperationStateProcessing ||
			lastOperation.State == gardener.LastOperationStatePending {
			msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", s.shoot.Name, lastOperation.State)
			m.log.Info(msg)

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonProcessing,
				"Unknown",
				"Shoot creation in progress")

			return stopWithRequeue()
		}

		// Runtime update is successful
		if lastOperation.State == gardener.LastOperationStateSucceeded {
			msg := fmt.Sprintf("Shoot %s successfully updated", s.shoot.Name)
			m.log.Info(msg)

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonProcessing,
				"True",
				"Shoot update completed")

			return stopWithNoRequeue()
		}

		// Runtime update is failed
		if lastOperation.State == gardener.LastOperationStateFailed {

			var reason ErrReason

			if len(s.shoot.Status.LastErrors) > 0 {
				reason = gardenerErrCodesToErrReason(s.shoot.Status.LastErrors...)
			}

			msg := fmt.Sprintf("error during cluster provisioning: reconcilation error for shoot %s, reason: %s, scheduling for retry", s.shoot.Name, reason)
			m.log.Info(msg)

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonProcessingErr,
				"Error",
				string(reason))

			return stopWithNoRequeue()
		}
	}

	return stopWithNoRequeue()
}

func gardenerErrCodesToErrReason(lastErrors ...gardener.LastError) ErrReason {
	var codes []gardener.ErrorCode
	var vals []string

	for _, e := range lastErrors {
		if len(e.Codes) > 0 {
			codes = append(codes, e.Codes...)
		}
	}

	for _, code := range codes {
		vals = append(vals, string(code))
	}

	return ErrReason(strings.Join(vals, ", "))
}
