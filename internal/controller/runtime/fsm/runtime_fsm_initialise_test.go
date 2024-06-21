package fsm

import (
	"context"
	"fmt"

	gardener_mocks "github.com/kyma-project/infrastructure-manager/internal/gardener/mocks"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"time"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
	"github.com/onsi/gomega/types"
)

var withTestFinalizer = withFinalizer("test-me-plz")

var _ = Describe("KIM sFnInitialise", func() {
	now := metav1.NewTime(time.Now())

	testCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// GIVEN

	testScheme := runtime.NewScheme()
	util.Must(imv1.AddToScheme(testScheme))

	withTestSchemeAndObjects := func(objs ...client.Object) fakeFSMOpt {
		return func(fsm *fsm) error {
			return withFakedK8sClient(testScheme, objs...)(fsm)
		}
	}

	testRt := imv1.Runtime{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
	}

	testRtWithFinalizer := imv1.Runtime{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-instance",
			Namespace:  "default",
			Finalizers: []string{"test-me-plz"},
		},
	}

	testRtWithDeletionTimestamp := imv1.Runtime{
		ObjectMeta: metav1.ObjectMeta{
			DeletionTimestamp: &now,
		},
	}

	testRtWithDeletionTimestampAndFinalizer := imv1.Runtime{
		ObjectMeta: metav1.ObjectMeta{
			DeletionTimestamp: &now,
			Finalizers:        []string{"test-me-plz"},
		},
	}

	testShoot := gardener.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-instance",
			Namespace: "default",
		},
	}

	testShootClient := gardener_mocks.ShootClient{}
	testShootClient.On("Get", mock.Anything, testShoot.Name, mock.Anything).Return(&testShoot, nil)

	testShootClientWithError := gardener_mocks.ShootClient{}
	testShootClientWithError.On("Get", mock.Anything, testShoot.Name, mock.Anything).Return(nil, fmt.Errorf("test error"))

	testFunction := buildTestFunction(sFnInitialize)

	// WHEN/THAN

	DescribeTable(
		"transition graph validation",
		testFunction,
		Entry(
			"should return nothing when CR is being deleted without finalizer",
			testCtx,
			must(newFakeFSM, withTestFinalizer),
			&systemState{instance: testRtWithDeletionTimestamp},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: BeNil(),
			},
		),
		Entry(
			"should return sFnDeleteShoot and no error when CR is being deleted with finalizer",
			testCtx,
			must(newFakeFSM, withTestFinalizer),
			&systemState{instance: testRtWithDeletionTimestampAndFinalizer},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnDeleteShoot"),
			},
		),
		Entry(
			"should return sFnUpdateStatus and no error when CR has been created",
			testCtx,
			must(newFakeFSM, withTestFinalizer, withTestSchemeAndObjects(&testRt)),
			&systemState{instance: testRt},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
				StateMatch:       []types.GomegaMatcher{haveFinalizer("test-me-plz")},
			},
		),
		Entry(
			"should return sFnCreateShoot and no error when CR has been created and shoot exists",
			testCtx,
			must(newFakeFSM, withTestFinalizer, withMockedShootClient(&testShootClient)),
			&systemState{instance: testRtWithFinalizer},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnCreateShoot"),
			},
		),
		Entry(
			"should return sFnCreateShoot and no error when shoot is missing",
			testCtx,
			must(newFakeFSM, withTestFinalizer, withMockedShootClient(&testShootClientWithError)),
			&systemState{instance: testRtWithFinalizer},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnCreateShoot"),
			},
		),
	)
})

type testOpts struct {
	MatchExpectedErr types.GomegaMatcher
	MatchNextFnState types.GomegaMatcher
	StateMatch       []types.GomegaMatcher
}

func buildTestFunction(fn stateFn) func(context.Context, *fsm, *systemState, testOpts) {
	return func(ctx context.Context, r *fsm, s *systemState, ops testOpts) {
		sFn, _, err := fn(ctx, r, s)

		Expect(err).To(ops.MatchExpectedErr)
		Expect(sFn).To(ops.MatchNextFnState)

		for _, match := range ops.StateMatch {
			Expect(&s.instance).Should(match)
		}
	}
}
