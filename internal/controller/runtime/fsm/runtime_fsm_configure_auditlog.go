package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/auditlogging"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureAuditLog(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure Audit Log state")

	shootNeedsToBeReconciled, err := m.AuditLogging.Enable(ctx, s.shoot)

	if err == nil && shootNeedsToBeReconciled {
		m.log.Info("Audit Log configured for shoot: " + s.shoot.Name)
		s.instance.UpdateStatePending(
			imv1.ConditionTypeAuditLogConfigured,
			imv1.ConditionReasonAuditLogConfigured,
			"Unknown",
			"Waiting for Gardener shoot to be Ready state after configuration of the Audit Logs",
		)

		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	if err == nil {
		s.instance.UpdateStateReady(
			imv1.ConditionTypeAuditLogConfigured,
			imv1.ConditionReasonAuditLogConfigured,
			"Audit Log state completed successfully",
		)

		return updateStatusAndStop()
	}

	setStateForAuditLogError := func(reason imv1.RuntimeConditionReason, pendingMsg string, readyMsg string) {
		if m.RCCfg.AuditLogMandatory {
			s.instance.UpdateStatePending(
				imv1.ConditionTypeAuditLogConfigured,
				reason,
				"False",
				pendingMsg,
			)
		} else {
			s.instance.UpdateStateReady(
				imv1.ConditionTypeAuditLogConfigured,
				reason,
				readyMsg)
		}
	}

	logError := func(err error, errorMsg, infoMsg string) {
		if m.RCCfg.AuditLogMandatory {
			m.log.Error(err, errorMsg)
		} else {
			m.log.Info(err.Error(), "Failed to configure Audit Log, but is not mandatory to be configured")
		}
	}

	if errors.Is(err, auditlogging.ErrMissingMapping) {
		pendingStatusMsg := err.Error()
		readyStatusMsg := "Missing region mapping for this shoot. Audit Log is not mandatory. Skipping configuration"
		setStateForAuditLogError(imv1.ConditionReasonAuditLogMissingRegionMapping, pendingStatusMsg, readyStatusMsg)

		errorMsg := "Failed to configure Audit Log, missing region mapping for this shoot"
		infoMsg := "Failed to configure Audit Log, missing region mapping for this shoot, but is not mandatory to be configured"
		logError(err, errorMsg, infoMsg)
	} else {
		pendingStatusMsg := err.Error()
		readyStatusMsg := "Configuration of Audit Log is not mandatory, error for context: " + err.Error()
		setStateForAuditLogError(imv1.ConditionReasonAuditLogError, pendingStatusMsg, readyStatusMsg)

		errorMsg := "Failed to configure Audit Log"
		infoMsg := "Failed to configure Audit Log, but is not mandatory to be configured"
		logError(err, errorMsg, infoMsg)
	}

	return updateStatusAndStop()
}
