package controller

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("nextRequeueAfter", func() {
	var (
		now = time.Now()
	)

	DescribeTable("should return expected values when",
		func(now, lastSyncTime time.Time, rotationPeriod time.Duration, modifier float64, expectedDuration time.Duration) {
			result := nextRequeue(now, lastSyncTime, rotationPeriod, modifier)
			Expect(result).To(Equal(expectedDuration))
		},
		Entry("receives all zero arguments", now, time.Time{}, time.Duration(0), 0.0, time.Duration(0)),
		Entry("receives arguments (now, zero, 1[m], 0.95)", now, time.Time{}, time.Minute, 0.95, time.Second*57),
		Entry("receives arguments (now, now-30[s], 1[m], 0.95)", now, now.Add(-30*time.Second), time.Minute, 0.95, time.Nanosecond*26999999999),
		Entry("receives arguments (now, now, 1[m], 0.95)", now, now, time.Minute, 0.95, time.Second*57),
	)
})
