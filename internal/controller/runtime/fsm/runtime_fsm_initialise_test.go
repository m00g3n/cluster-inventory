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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("KIM sFnInitialise", func() {
	now := metav1.NewTime(time.Now())

	testContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// GIVEN

	testScheme := runtime.NewScheme()
	util.Must(imv1.AddToScheme(testScheme))

	withTestFinalizer := withFinalizer("test-me-plz")
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

	// THAN

	DescribeTable(
		"transition graph validation",
		testInitialise,
		Entry(
			"should return nothing when CR is being deleted without finalizer",
			testContext,
			must(newFakeFSM, withTestFinalizer),
			&systemState{instance: testRtWithDeletionTimestamp},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: BeNil(),
			},
		),
		Entry(
			"should return sFnDeleteShoot and no error when CR is being deleted with finalizer",
			testContext,
			must(newFakeFSM, withTestFinalizer),
			&systemState{instance: testRtWithDeletionTimestampAndFinalizer},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnDeleteShoot"),
			},
		),
		Entry(
			"should return sFnUpdateStatus and no error when CR has been created",
			testContext,
			must(newFakeFSM, withTestFinalizer, withTestSchemeAndObjects(&testRt)),
			&systemState{instance: testRt},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
				StateMatch:       []types.GomegaMatcher{haveFinalizer("test-me-plz")},
			},
		),
		Entry(
			"should return sFnPrepareCluster and no error when CR has been created and shoot exists",
			testContext,
			must(newFakeFSM, withTestFinalizer, withMockedShootClient(&testShootClient)),
			&systemState{instance: testRtWithFinalizer},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnPrepareCluster"),
			},
		),
		Entry(
			"should return sFnUpdateStatus and no error when shoot is missing",
			testContext,
			must(newFakeFSM, withTestFinalizer, withMockedShootClient(&testShootClientWithError)),
			&systemState{instance: testRtWithFinalizer},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
	)
})

type testInitialiseOpts struct {
	MatchExpectedErr types.GomegaMatcher
	MatchNextFnState types.GomegaMatcher
	StateMatch       []types.GomegaMatcher
}

func testInitialise(ctx context.Context, r *fsm, s *systemState, ops testInitialiseOpts) {
	sFn, _, err := sFnInitialize(ctx, r, s)

	Expect(err).To(ops.MatchExpectedErr)
	Expect(sFn).To(ops.MatchNextFnState)

	for _, match := range ops.StateMatch {
		Expect(&s.instance).Should(match)
	}
}
