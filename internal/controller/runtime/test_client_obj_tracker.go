package runtime

import (
	"fmt"
	"sync"

	gardener_api "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clienttesting "k8s.io/client-go/testing"
)

// CustomTracker implements ObjectTracker with a sequence of Shoot objects
// it will be updated with a different shoot sequence for each test case
type CustomTracker struct {
	clienttesting.ObjectTracker
	shootSequence []*gardener_api.Shoot
	seedSequence  []*gardener_api.Seed
	shootCallCnt  int
	seedCallCnt   int
	mu            sync.Mutex
}

func NewCustomTracker(tracker clienttesting.ObjectTracker, shoots []*gardener_api.Shoot, seeds []*gardener_api.Seed) *CustomTracker {
	return &CustomTracker{
		ObjectTracker: tracker,
		shootSequence: shoots,
		seedSequence:  seeds,
	}
}

func (t *CustomTracker) IsSequenceFullyUsed() bool {
	return (t.shootCallCnt == len(t.shootSequence) && len(t.shootSequence) > 0) && (t.seedCallCnt == len(t.seedSequence) && len(t.seedSequence) > 0)
}

func (t *CustomTracker) Get(gvr schema.GroupVersionResource, ns, name string) (runtime.Object, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if gvr.Resource == "shoots" {
		return getNextShoot(t.shootSequence, t)
	} else if gvr.Resource == "seeds" {
		return getNextSeed(t.seedSequence, t)
	}

	return t.ObjectTracker.Get(gvr, ns, name)
}

func getNextSeed(seed []*gardener_api.Seed, t *CustomTracker) (runtime.Object, error) {
	if t.seedCallCnt < len(seed) {
		obj := seed[t.seedCallCnt]
		t.seedCallCnt++

		if obj == nil {
			return nil, k8serrors.NewNotFound(schema.GroupResource{}, "")
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no more seeds in sequence")
}

func getNextShoot(shoot []*gardener_api.Shoot, t *CustomTracker) (runtime.Object, error) {
	if t.shootCallCnt < len(shoot) {
		obj := shoot[t.shootCallCnt]
		t.shootCallCnt++

		if obj == nil {
			return nil, k8serrors.NewNotFound(schema.GroupResource{}, "")
		}
		return obj, nil
	}
	return nil, fmt.Errorf("no more shoots in sequence")
}

//func getNextObject[T any](sequence []*T, counter *int) (*T, error) {
//	if *counter < len(sequence) {
//		obj := sequence[*counter]
//		*counter++
//
//		if obj == nil {
//			return nil, k8serrors.NewNotFound(schema.GroupResource{}, "")
//		}
//		return obj, nil
//	}
//	return nil, fmt.Errorf("no more objects in sequence")
//}

func (t *CustomTracker) Delete(gvr schema.GroupVersionResource, ns, name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if gvr.Resource == "shoots" {
		for index, shoot := range t.shootSequence {
			if shoot != nil && shoot.Name == name {
				t.shootSequence[index] = nil
				return nil
			}
		}
		return k8serrors.NewNotFound(schema.GroupResource{}, "")
	}
	return t.ObjectTracker.Delete(gvr, ns, name)
}
