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
		// Conditions to change
		s.instance.UpdateStatePending(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerError,
			"False",
			err.Error(),
		)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	if auditLogConfigWasChanged {
		err = m.ShootClient.Update(ctx, s.shoot)
		if err != nil {
			m.log.Error(err, "Failed to update shoot")

			// Conditions to change
			s.instance.UpdateStatePending(
				imv1.ConditionTypeRuntimeProvisioned,
				imv1.ConditionReasonGardenerError,
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

	// config audit logow sie nie zmienil nie trzeba aktualizowac shoota
	m.log.Info("Audit Log Tenant did not change, skipping update of cluster")
	s.instance.UpdateStateReady(
		imv1.ConditionTypeRuntimeConfigured,
		imv1.ConditionReasonConfigurationCompleted,
		"Audit Log configured",
	)

	return updateStatusAndStop()
}
