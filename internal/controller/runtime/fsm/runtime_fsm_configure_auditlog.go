package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureAuditLog(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure Audit Log state")

	auditLogConfigWasChanged, err := m.AuditLogging.Enable(ctx, s.shoot)
	if err != nil {
		m.log.Error(err, "Failed to configure Audit Log")
		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeConfigured,
			imv1.ConditionReasonAuditLogConfigured,
			"False",
			err.Error(),
		)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	if auditLogConfigWasChanged {
		err = m.ShootClient.Update(ctx, s.shoot)
		if err != nil {
			m.log.Error(err, "Failed to update shoot")

			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeConfigured,
				imv1.ConditionReasonAuditLogConfigured,
				"False",
				err.Error(),
			)
			return updateStatusAndRequeueAfter(gardenerRequeueDuration)
		}
		m.log.Info("Audit Log configured for shoot:" + s.shoot.Name)
		s.instance.UpdateStateReady(
			imv1.ConditionTypeRuntimeConfigured,
			imv1.ConditionReasonConfigurationCompleted,
			"Audit Log configured",
		)
		return updateStatusAndStop()
	}

	m.log.Info("Audit Log Tenant did not change, skipping update of cluster")
	s.instance.UpdateStateReady(
		imv1.ConditionTypeRuntimeConfigured,
		imv1.ConditionReasonConfigurationCompleted,
		"Audit Log configured",
	)

	return updateStatusAndStop()
}
