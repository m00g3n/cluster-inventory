package fsm

/*
import (
	"context"
	"time"

	"github.com/kyma-project/infrastructure-manager/internal/controller/runtime/fsm/testing"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

var _ = Describe("KIM sFnInitialise", func() {

	testCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testFunction := buildTestFunction(sFnPrepareCluster)

	DescribeTable(
		"transition graph validation",
		testFunction,
		//		Entry(
		//			"validate missing shoot",
		//			testCtx,
		//			must(newFakeFSM),
		//			&systemState{},
		//			testOpts{
		//				MatchExpectedErr: HaveOccurred(),
		//				MatchNextFnState: BeNil(),
		//			},
		//		),
		Entry(
			"validate missing DNS",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootNoDNS,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"retry when last operation on shoot is unknown",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootNoDNS,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"validate missing DNS domain",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootNoDNSDomain,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"validate missing last operation",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootMissingLastOperation,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation processing",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationProcessing,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation create pending",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationPending,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation create succeeded",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationSucceeded,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation create failed",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationFailed,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation reconcile processing",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationReconcileProcessing,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation reconcile pending",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationReconcilePending,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation reconcile succeeded",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationReconcileSucceeded,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
		Entry(
			"last operation reconcile failed",
			testCtx,
			must(newFakeFSM),
			&systemState{
				shoot: &testing.ShootLastOperationReconcileFailed,
			},
			testOpts{
				MatchExpectedErr: BeNil(),
				MatchNextFnState: haveName("sFnUpdateStatus"),
			},
		),
	)
})
*/
