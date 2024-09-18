package fsm

import (
	"context"
	"strconv"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/infrastructure-manager/internal/auditlogging"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureAuditLog(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure Audit Log state")

	wasAuditLogEnabled, err := m.AuditLogging.Enable(ctx, s.shoot)

	if wasAuditLogEnabled && err == nil {
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

	return handleError(err, m, s)
	//if err != nil { //nolint:nestif
	//	if k8serrors.IsConflict(err) {
	//		m.log.Error(err, "Conflict while updating Shoot object after applying Audit Log configuration, retrying")
	//		s.instance.UpdateStatePending(
	//			imv1.ConditionTypeAuditLogConfigured,
	//			imv1.ConditionReasonAuditLogError,
	//			"True",
	//			err.Error(),
	//		)
	//		return updateStatusAndRequeue()
	//	}
	//	errorMessage := err.Error()
	//
	//	if errors.Is(err, auditlogging.ErrMissingMapping) {
	//		if m.RCCfg.AuditLogMandatory {
	//			m.log.Error(err, "AuditLogMandatory", auditLogMandatoryString, "providerType", s.shoot.Spec.Provider.Type, "region", s.shoot.Spec.Region)
	//			s.instance.UpdateStatePending(
	//				imv1.ConditionTypeAuditLogConfigured,
	//				imv1.ConditionReasonAuditLogMissingRegionMapping,
	//				"False",
	//				errorMessage,
	//			)
	//		} else {
	//			m.log.Info(errorMessage, "AuditLogMandatory", auditLogMandatoryString, "providerType", s.shoot.Spec.Provider.Type, "region", s.shoot.Spec.Region)
	//			s.instance.UpdateStateReady(
	//				imv1.ConditionTypeAuditLogConfigured,
	//				imv1.ConditionReasonAuditLogMissingRegionMapping,
	//				"Missing region mapping for this shoot. Audit Log is not mandatory. Skipping configuration")
	//		}
	//	} else {
	//		if m.RCCfg.AuditLogMandatory {
	//			m.log.Error(err, "AuditLogMandatory", auditLogMandatoryString)
	//			s.instance.UpdateStatePending(
	//				imv1.ConditionTypeAuditLogConfigured,
	//				imv1.ConditionReasonAuditLogError,
	//				"False",
	//				errorMessage)
	//		} else {
	//			m.log.Info(errorMessage, "AuditLogMandatory", auditLogMandatoryString)
	//			s.instance.UpdateStateReady(
	//				imv1.ConditionTypeAuditLogConfigured,
	//				imv1.ConditionReasonAuditLogError,
	//				"Configuration of Audit Log is not mandatory, error for context: "+errorMessage)
	//		}
	//	}
	//}
}

func handleError(err error, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
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

	logError := func(err error, errorMsg, infoMsg string, keysAndValues ...any) {
		if m.RCCfg.AuditLogMandatory {
			m.log.Error(err, errorMsg, keysAndValues)
		} else {
			m.log.Info(err.Error(), infoMsg)
		}
	}

	if k8serrors.IsConflict(err) {
		m.log.Error(err, "Conflict while updating Shoot object after applying Audit Log configuration, retrying")
		s.instance.UpdateStatePending(
			imv1.ConditionTypeAuditLogConfigured,
			imv1.ConditionReasonAuditLogError,
			"True",
			err.Error(),
		)

		return updateStatusAndRequeue()
	}

	auditLogMandatoryString := strconv.FormatBool(m.RCCfg.AuditLogMandatory)

	if errors.Is(err, auditlogging.ErrMissingMapping) {
		pendingStatusMsg := err.Error()
		readyStatusMsg := "Missing region mapping for this shoot. Audit Log is not mandatory. Skipping configuration"
		setStateForAuditLogError(imv1.ConditionReasonAuditLogMissingRegionMapping, pendingStatusMsg, readyStatusMsg)

		errorMsg := "Failed to configure Audit Log"
		infoMsg := "Failed to configure Audit Log"
		logError(err, errorMsg, infoMsg, "AuditLogMandatory", auditLogMandatoryString, "providerType", s.shoot.Spec.Provider.Type, "region", s.shoot.Spec.Region)

		return updateStatusAndRequeue()
	}

	pendingStatusMsg := err.Error()
	readyStatusMsg := "Configuration of Audit Log is not mandatory, error for context: " + err.Error()
	setStateForAuditLogError(imv1.ConditionReasonAuditLogError, pendingStatusMsg, readyStatusMsg)

	errorMsg := "Failed to configure Audit Log"
	infoMsg := "Failed to configure Audit Log"
	logError(err, errorMsg, infoMsg, "AuditLogMandatory", auditLogMandatoryString)

	return updateStatusAndRequeue()
}
