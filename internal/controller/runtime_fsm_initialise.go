package controller

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// There is a decision made to not rely on state of the Runtime CR we have already set
// All the states we set in the operator are about to be read only by the external clients

func sFnInitialize(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	instanceIsBeingDeleted := !s.instance.GetDeletionTimestamp().IsZero()
	instanceHasFinalizer := controllerutil.ContainsFinalizer(&s.instance, m.Finalizer)

	// in case instance does not have finalizer - add it and update instance
	if !instanceIsBeingDeleted && !instanceHasFinalizer {
		m.log.Info("adding finalizer")
		controllerutil.AddFinalizer(&s.instance, m.Finalizer)

		err := m.Update(ctx, &s.instance)
		if err != nil {
			return stopWithErrorAndRequeue(err)
		}

		s.instance.UpdateStateProcessing(
			imv1.ConditionTypeRuntimeProvisioning,
			imv1.ConditionReasonInitialized,
			"initialized",
		)
		return stopWithRequeue()
	}

	// in case instance has no finalizer and instance is being deleted - end reconciliation
	if instanceIsBeingDeleted {
		// in case instance is being deleted and does not have finalizer - delete shoot
		if controllerutil.ContainsFinalizer(&s.instance, m.Finalizer) {
			return switchState(sFnDeleteShoot)
		}
		m.log.Info("Instance is being deleted")
		// stop state machine ???
		return nil, nil, nil
	}

	_, err := m.shootClient.Get(ctx, s.instance.Name, v1.GetOptions{})

	if err != nil {
		if apierrors.IsNotFound(err) {
			m.log.Info("Gardener shoot does not exist, creating new one")
			return switchState(sFnCreateShoot)
		}

		m.log.Info("Failed to get shoot", "error", err)
		return stopWithRequeue()
	}

	return switchState(sFnPrepareCluster)
}
