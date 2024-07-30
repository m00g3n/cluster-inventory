package fsm

import (
	"context"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnDeleteKubeconfig(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	m.log.Info("delete Kubeconfig/GardenerCluster CR state")

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

	runtimeID := s.instance.Labels[imv1.LabelKymaRuntimeID]

	var cluster imv1.GardenerCluster
	err := m.Get(ctx, types.NamespacedName{
		Namespace: s.instance.Namespace,
		Name:      runtimeID,
	}, &cluster)

	if err == nil {
		m.log.Info("deleting GardenerCluster CR for %s", s.shoot.Name)
		err := m.Delete(ctx, &cluster)
		if err != nil {
			m.log.Error(err, "Failed to delete gardener Cluster CR")
		}

		return requeueAfter(controlPlaneRequeueDuration) // requeue to wait until GardenerCluster CR is deleted
	}

	if !k8serrors.IsNotFound(err) {
		m.log.Error(err, "GardenerCluster CR read error", "name", runtimeID)
		s.instance.UpdateStateDeletion(imv1.RuntimeStateTerminating, imv1.ConditionReasonKubernetesAPIErr, "False", err.Error())
		return updateStatusAndStop()
	}

	if !s.shoot.GetDeletionTimestamp().IsZero() {
		m.log.Info("Waiting for shoot to be deleted", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)
		return requeueAfter(gardenerRequeueDuration)
	}

	m.log.Info("deleting shoot", "Name", s.shoot.Name, "Namespace", s.shoot.Namespace)
	err = m.ShootClient.Delete(ctx, s.shoot)
	if err != nil {
		m.log.Error(err, "Failed to delete gardener Shoot")

		s.instance.UpdateStateDeletion(
			imv1.ConditionTypeRuntimeProvisioned,
			imv1.ConditionReasonGardenerError,
			"False",
			"Gardener API shoot delete error",
		)
		return updateStatusAndRequeueAfter(gardenerRequeueDuration)
	}
	return requeueAfter(gardenerRequeueDuration)
}
