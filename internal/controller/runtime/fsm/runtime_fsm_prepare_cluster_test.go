package fsm

import (
	"context"

	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("KIM sFnInitialise", func() {
	_ = metav1.NewTime(time.Now())

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_ = buildTestFunction(sFnPrepareCluster)
})

func testPrepareCluster(ctx context.Context, r *fsm, s *systemState, ops testOpts) {
	sFn, _, err := sFnPrepareCluster(ctx, r, s)

	Expect(err).To(ops.MatchExpectedErr)
	Expect(sFn).To(ops.MatchNextFnState)

	for _, match := range ops.StateMatch {
		Expect(&s.instance).Should(match)
	}
}
