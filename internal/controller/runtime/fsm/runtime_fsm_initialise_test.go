package fsm

import (
	"context"
	"fmt"
	gardener_mocks "github.com/kyma-project/infrastructure-manager/internal/gardener/mocks"

	"time"

	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func newTestFSM(finalizer string, objs ...client.Object) *fsm {
	// create a new scheme for the test
	scheme := runtime.NewScheme()
	// add supported types to the scheme
	util.Must(imv1.AddToScheme(scheme))

	c := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		Build()

	return &fsm{
		K8s: K8s{
			ShootClient: nil,
			Client:      c,
		},
		RCCfg: RCCfg{
			Finalizer: finalizer,
		},
	}
}

func newTestFSMWithGardener(finalizer string, shoot gardener.Shoot) *fsm {
	// create a new scheme for the test
	scheme := runtime.NewScheme()
	// add supported types to the scheme
	util.Must(imv1.AddToScheme(scheme))

	c := gardener_mocks.ShootClient{}
	c.On("Get", mock.Anything, shoot.Name, mock.Anything).Return(&shoot, nil)

	return &fsm{
		K8s: K8s{
			ShootClient: &c,
		},
		RCCfg: RCCfg{
			Finalizer: finalizer,
		},
	}
}

func newTestFSMWithShootNotFound(finalizer string) *fsm {
	// create a new scheme for the test
	scheme := runtime.NewScheme()
	// add supported types to the scheme
	util.Must(imv1.AddToScheme(scheme))

	c := gardener_mocks.ShootClient{}
	c.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("not found"))

	return &fsm{
		K8s: K8s{
			ShootClient: &c,
		},
		RCCfg: RCCfg{
			Finalizer: finalizer,
		},
	}
}

var _ = Describe("KIM sFnInitialise", func() {
	now := metav1.NewTime(time.Now())

	testContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	DescribeTable(
		"transition graph validation",
		testInitialise,
		Entry(
			"should return nothing when CR is being deleted without finalizer",
			testContext,
			newTestFSM("test-me-plz"),
			&systemState{
				instance: imv1.Runtime{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: &now,
					},
				},
			},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: BeNil(),
			},
		),
		Entry(
			"should return sFnDeleteShoot and no error when CR is being deleted with finalizer",
			testContext,
			newTestFSM("test-me-plz"),
			&systemState{
				instance: imv1.Runtime{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: &now,
						Finalizers:        []string{"test-me-plz"},
					},
				},
			},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnDeleteShoot"),
			},
		),
		Entry(
			"should return sFnUpdateStatus and no error when CR has been created",
			testContext,
			newTestFSM("test-me-plz", &imv1.Runtime{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-instance",
					Namespace: "default",
				},
			}),
			&systemState{
				instance: imv1.Runtime{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-instance",
						Namespace: "default",
					},
				},
			},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
				StateMatch: []types.GomegaMatcher{
					haveFinalizer("test-me-plz"),
				},
			},
		),
		Entry(
			"should return sFnPrepareCluster and no error when CR has been created and shoot exists",
			testContext,
			newTestFSMWithGardener("test-me-plz", gardener.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-instance",
					Namespace: "default",
				},
			}),
			&systemState{
				instance: imv1.Runtime{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-instance",
						Namespace:  "default",
						Finalizers: []string{"test-me-plz"},
					},
				},
			},
			testInitialiseOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnPrepareCluster"),
			},
		),
		Entry(
			"should return sFnUpdateStatus and no error when shoot is missing",
			testContext,
			newTestFSMWithShootNotFound("test-me-plz"),
			&systemState{
				instance: imv1.Runtime{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-instance",
						Namespace:  "default",
						Finalizers: []string{"test-me-plz"},
					},
				},
			},
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
