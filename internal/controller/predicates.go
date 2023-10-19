package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

func newRotationNotNeededPredicate(rotationPeriod time.Duration) predicate.Predicate {
	return RotationNotNeededPredicate{
		rotationPeriod: rotationPeriod,
	}
}

type RotationNotNeededPredicate struct {
	rotationPeriod time.Duration
}

func (RotationNotNeededPredicate) Create(event.CreateEvent) bool {
	return true
}

func (r RotationNotNeededPredicate) Update(e event.UpdateEvent) bool {
	return kubeconfigMustBeRotated(e.ObjectNew, r.rotationPeriod)
}

func (RotationNotNeededPredicate) Delete(event.DeleteEvent) bool {
	return true
}

func (r RotationNotNeededPredicate) Generic(e event.GenericEvent) bool {
	return kubeconfigMustBeRotated(e.Object, r.rotationPeriod)
}

func kubeconfigMustBeRotated(gardenerCluster client.Object, validityTime time.Duration) bool {
	return rotationForced(gardenerCluster) || rotationTimePassed(gardenerCluster, validityTime)
}

func rotationTimePassed(gardenerCluster client.Object, rotationPeriod time.Duration) bool {
	annotations := gardenerCluster.GetAnnotations()

	_, found := annotations[lastKubeconfigSyncAnnotation]

	if !found {
		return true
	}

	lastSyncTimeString := annotations[lastKubeconfigSyncAnnotation]
	lastSyncTime, err := time.Parse(time.RFC3339, lastSyncTimeString)
	if err != nil {
		return true
	}
	now := time.Now()
	alreadyValidFor := lastSyncTime.Sub(now)

	return alreadyValidFor.Minutes() >= rotationPeriod.Minutes()
}

func rotationForced(gardenerCluster client.Object) bool {
	_, found := gardenerCluster.GetAnnotations()[forceKubeconfigRotationAnnotation]

	return found
}
