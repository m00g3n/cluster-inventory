package fsm

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// There is a decision made to not rely on state of the Runtime CR we have already set
// All the states we set in the operator are about to be read only by the external clients

func sFnInitialize(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {

	instanceIsNotBeingDeleted := s.instance.GetDeletionTimestamp().IsZero()
	instanceHasFinalizer := controllerutil.ContainsFinalizer(&s.instance, m.Finalizer)

	// in case instance does not have finalizer - add it and update instance
	if instanceIsNotBeingDeleted {
		if !instanceHasFinalizer {
			m.log.Info("adding finalizer")
			controllerutil.AddFinalizer(&s.instance, m.Finalizer)

			err := m.Update(ctx, &s.instance)
			if err != nil {
				return stopWithErrorAndNoRequeue(err)
			}
			return stopWithRequeueAfter(time.Millisecond * 250)
		}

		if s.shoot == nil {
			provisioningCondition := meta.FindStatusCondition(
				s.instance.Status.Conditions,
				string(imv1.ConditionTypeRuntimeProvisioned),
			)

			if provisioningCondition == nil {
				s.instance.UpdateStatePending(
					imv1.ConditionTypeRuntimeProvisioned,
					imv1.ConditionReasonInitialized,
					"Unknown",
					"Runtime initialized",
				)
				return stopWithRequeueAfter(time.Millisecond * 250)
			}

			m.log.Info("Gardener shoot does not exist, creating new one")
			return switchState(sFnCreateShoot)
		}

		return switchState(sFnPrepareCluster) // wait for operation to complete

	} else {
		if instanceHasFinalizer {
			//if s.shoot.GetDeletionTimestamp().IsZero() {
			m.log.Info("Instance is being deleted")
			return switchState(sFnDeleteShoot)
			//}

			//return switchState(sFnWaitForShootDeletion)
		}

		//if s.shoot == nil {
		//	m.log.Info("Removing finalizer")
		//	controllerutil.RemoveFinalizer(&s.instance, m.Finalizer)
		//	return stopWithNoRequeue()
		//}
	}
	// stop state machine
	return nil, nil, nil
}
