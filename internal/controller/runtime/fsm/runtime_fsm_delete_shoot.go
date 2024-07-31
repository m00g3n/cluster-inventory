package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func sFnDeleteShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("delete shoot state")
	if !isGardenerCloudDelConfirmationSet(s.shoot.Annotations) {
		m.log.Info("patching shoot with del-confirmation")
		// workaround for Gardener client
		setObjectFields(s.shoot)
		s.shoot.Annotations = addGardenerCloudDelConfirmation(s.shoot.Annotations)

		err := m.ShootClient.Patch(ctx, s.shoot, client.Apply, &client.PatchOptions{
			FieldManager: "kim",
			Force:        ptr.To(true),
		})

		if err != nil {
			m.log.Error(err, "unable to patch shoot:", s.shoot.Name)
			return requeue()
		}
	}

	if !s.instance.IsStateWithConditionSet(imv1.RuntimeStateTerminating, imv1.ConditionTypeRuntimeProvisioned, imv1.ConditionReasonDeletion) {
		m.log.Info("setting state to in deletion")
		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonDeletion,
			"Unknown",
			"Runtime deletion initialised",
		)
		return updateStatusAndRequeue()
	}

	m.log.Info("deleting shoot")
	err := m.ShootClient.Delete(ctx, s.shoot)
	if err != nil {
		m.log.Error(err, "Failed to delete gardener Shoot")

		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerError,
			"False",
			"Gardener API delete error",
		)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}

	return updateStatusAndRequeueAfter(gardenerRequeueDuration)
}

func isGardenerCloudDelConfirmationSet(a map[string]string) bool {
	if len(a) == 0 {
		return false
	}
	val, found := a[imv1.AnnotationGardenerCloudDelConfirmation]
	return found && (val == "true")
}

func addGardenerCloudDelConfirmation(a map[string]string) map[string]string {
	if len(a) == 0 {
		a = map[string]string{}
	}
	a[imv1.AnnotationGardenerCloudDelConfirmation] = "true"
	return a
}
