package fsm

import (
	"context"
	"fmt"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardenerhelper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

type ErrReason string

func sFnPrepareCluster(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	shoot, err := m.ShootClient.Get(ctx, s.instance.Name, v1.GetOptions{})

	if err != nil || shoot == nil {
		msg := fmt.Sprintf("Cannot get shoot: %s, scheduling for retry", s.instance.Name)
		m.log.Info(msg)
		return stopWithRequeue()
	}

	if shoot.Spec.DNS == nil || shoot.Spec.DNS.Domain == nil {
		msg := fmt.Sprintf("DNS Domain is not set yet for shoot: %s, scheduling for retry", shoot.Name)
		m.log.Info(msg)
		return stopWithRequeue()
	}

	lastOperation := shoot.Status.LastOperation
	if lastOperation == nil {
		msg := fmt.Sprintf("Last operation is nil for shoot: %s, scheduling for retry", shoot.Name)
		m.log.Info(msg)
		return stopWithRequeue()
	}

	if lastOperation.Type == gardener.LastOperationTypeCreate {

		if lastOperation.State == gardener.LastOperationStateProcessing ||
			lastOperation.State == gardener.LastOperationStatePending {
			msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", shoot.Name, lastOperation.State)
			m.log.Info(msg)

			s.instance.UpdateStateCreating(
				imv1.ConditionTypeRuntimeProvisioning,
				imv1.ConditionReasonShootCreationPending,
				"Shoot creation in progress")

			return stopWithRequeue()
		}

		if lastOperation.State == gardener.LastOperationStateSucceeded {
			msg := fmt.Sprintf("Shoot %s successfully created", shoot.Name)
			m.log.Info(msg)

			if !s.instance.IsRuntimeStateSet(imv1.RuntimeStateCreating, imv1.ConditionTypeRuntimeProvisioning, imv1.ConditionReasonShootCreationCompleted) {
				s.instance.UpdateStateCreating(
					imv1.ConditionTypeRuntimeProvisioning,
					imv1.ConditionReasonShootCreationCompleted,
					"Shoot creation completed successfully")

				return stopWithRequeue()
			}

			return switchState(sFnProcessShoot)
		}

		if lastOperation.State == gardener.LastOperationStateFailed {
			if gardenerhelper.HasErrorCode(shoot.Status.LastErrors, gardener.ErrorInfraRateLimitsExceeded) {
				msg := fmt.Sprintf("Error during cluster provisioning: Rate limits exceeded for Shoot %s, scheduling for retry", shoot.Name)
				m.log.Info(msg)
				return stopWithRequeue()
			}

			msg := fmt.Sprintf("Provisioning failed for shoot: %s ! Last state: %s, Description: %s", shoot.Name, lastOperation.State, lastOperation.Description)
			m.log.Info(msg)

			s.instance.UpdateStateError(
				imv1.ConditionTypeRuntimeProvisioning,
				imv1.ConditionReasonCreationError,
				"Shoot creation failed")

			return stopWithNoRequeue()
		}
	}

	// Runtime update is in progress
	if lastOperation.Type == gardener.LastOperationTypeReconcile {

		if lastOperation.State == gardener.LastOperationStateProcessing ||
			lastOperation.State == gardener.LastOperationStatePending {
			msg := fmt.Sprintf("Shoot %s is in %s state, scheduling for retry", shoot.Name, lastOperation.State)
			m.log.Info(msg)

			s.instance.UpdateStateProcessing(
				imv1.ConditionTypeRuntimeUpdate,
				imv1.ConditionReasonProcessing,
				"Shoot creation in progress")

			return stopWithRequeue()
		}

		if lastOperation.State == gardener.LastOperationStateSucceeded {
			msg := fmt.Sprintf("Shoot %s successfully updated", shoot.Name)
			m.log.Info(msg)

			s.instance.UpdateStateProcessing(
				imv1.ConditionTypeRuntimeUpdate,
				imv1.ConditionReasonProcessingCompleted,
				"Shoot creation in progress")

			return stopWithNoRequeue()
		}

		if lastOperation.State == gardener.LastOperationStateFailed {

			var reason ErrReason

			if len(shoot.Status.LastErrors) > 0 {
				reason = gardenerErrCodesToErrReason(shoot.Status.LastErrors...)
			}

			msg := fmt.Sprintf("error during cluster provisioning: reconcilation error for shoot %s, reason: %s, scheduling for retry", shoot.Name, reason)
			m.log.Info(msg)

			s.instance.UpdateStateError(
				imv1.ConditionTypeRuntimeUpdate,
				imv1.ConditionReasonProcessingErr,
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
