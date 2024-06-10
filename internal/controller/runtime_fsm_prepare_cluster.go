package controller

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

	shoot, err := m.shootClient.Get(ctx, s.instance.Name, v1.GetOptions{})

	if err != nil || shoot == nil {
		msg := fmt.Sprintf("Cannot get shoot: %s, scheduling for retry", shoot.Name)
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

			//s.instance.UpdateStateCreating(
			//	imv1.ConditionTypeRuntimeProvisioning,
			//	imv1.ConditionReasonShootCreationPending,
			//	"Shoot creation in progress")

			return stopWithRequeue()
		}

		if lastOperation.State == gardener.LastOperationStateSucceeded {
			msg := fmt.Sprintf("Shoot %s successfully created", shoot.Name)
			m.log.Info(msg)

			if s.instance.IsStateCreating() {
				return switchState(sFnProcessShoot)
			}

			s.instance.UpdateStateCreating(
				imv1.ConditionTypeRuntimeProvisioning,
				imv1.ConditionReasonShootCreationCompleted,
				"Shoot creation completed successfully")

			return stopWithRequeue()
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

	// updating is progressing
	if lastOperation.Type == gardener.LastOperationTypeReconcile {

		var reason ErrReason

		if len(shoot.Status.LastErrors) > 0 {
			reason = gardenerErrCodesToErrReason(shoot.Status.LastErrors...)
		}

		msg := fmt.Sprintf("error during cluster provisioning: reconcilation error for shoot %s, reason: %s, scheduling for retry", shoot.Name, reason)
		m.log.Info(msg)
		return stopWithRequeue()
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

//func isProcessingCluster(s *systemState) bool {
//	condition := meta.FindStatusCondition(s.instance.Status.Conditions, string(imv1.ConditionTypeRuntimeProvisioning))
//	if condition == nil {
//		return false
//	}
//
//	if condition.Reason != string(imv1.ConditionReasonProcessing) &&
//		condition.Reason != string(imv1.ConditionReasonProcessingErr) {
//		return false
//	}
//
//	return true
//}
