package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func sFnDeleteShoot(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("delete shoot state")

	// init section
	//if !s.instance.IsStateWithConditionSet(imv1.RuntimeStateTerminating, imv1.ConditionTypeRuntimeDeprovisioned, imv1.ConditionReasonGardenerShootDeleted) {
	//	m.log.Info("setting state to in deletion")
	//	s.instance.UpdateStateDeletion(
	//		imv1.ConditionTypeRuntimeDeprovisioned,
	//		imv1.ConditionReasonGardenerShootDeleted,
	//		"Unknown",
	//		"Runtime shoot deletion started",
	//	)
	//	return updateStatusAndRequeue()
	//}

	// wait section
	if !s.shoot.GetDeletionTimestamp().IsZero() {
		m.log.Info("Waiting for shoot to be deleted", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)
		return requeueAfter(gardenerRequeueDuration)
	}

	// action section
	if !isGardenerCloudDelConfirmationSet(s.shoot.Annotations) {
		m.log.Info("patching shoot with del-confirmation")
		// workaround for Gardener client
		setObjectFields(s.shoot)
		s.shoot.Annotations = addGardenerCloudDelConfirmation(s.shoot.Annotations)

		err := m.ShootClient.Patch(ctx, s.shoot, client.Apply, &client.PatchOptions{
			FieldManager: "kim",
			Force:        ptrTo(true),
		})

		if err != nil {
			m.log.Error(err, "unable to patch shoot:", s.shoot.Name)
			return requeue()
		}
	}

	m.log.Info("deleting shoot", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)
	err := m.ShootClient.Delete(ctx, s.shoot)
	if err != nil {
		// action error handler section
		m.log.Error(err, "Failed to delete gardener Shoot")
		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerShootDeleted,
			"False",
			"Gardener API shoot delete error",
		)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}
	// out succeeded section
	return requeueAfter(gardenerRequeueDuration)
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
