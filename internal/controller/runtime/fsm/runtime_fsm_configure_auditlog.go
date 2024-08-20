package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnConfigureAuditLog(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("Configure Audit Log state")

	err := m.AuditLogging.Enable(ctx, s.shoot)
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

	m.log.Info("Audit Log configured for shoot: " + s.shoot.Name)
	s.instance.UpdateStateReady(
		imv1.ConditionTypeRuntimeConfigured,
		imv1.ConditionReasonConfigurationCompleted,
		"Audit Log configured",
	)

	return updateStatusAndStop()
}
