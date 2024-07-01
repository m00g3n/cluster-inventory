package kubeconfig

import (
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

var _ = Describe("nextRequeueAfter", func() {
	var (
		now        = time.Now()
		newYear, _ = time.Parse(time.RFC3339, "2024-01-01T00:00:01Z00:00")
	)

	DescribeTable("should return expected values when",
		func(now, lastSyncTime time.Time, rotationPeriod time.Duration, modifier float64, expectedDuration time.Duration) {
			result := nextRequeue(now, lastSyncTime, rotationPeriod, modifier)
			Expect(result).To(BeNumerically("~", expectedDuration, 1))
		},
		Entry("receives all zero arguments", now, time.Time{}, time.Duration(0), 0.0, time.Duration(0)),
		Entry("receives arguments (now, zero, 1[m], 0.95)", now, time.Time{}, time.Minute, 0.95, time.Second*57),
		Entry("receives arguments (now, now-30[s], 1[m], 0.95)", now, now.Add(-30*time.Second), time.Minute, 0.95, time.Nanosecond*26999999999),
		Entry("receives arguments (now, now-900[s], 1[m], 0.95)", now, now.Add(-900*time.Second), time.Minute, 0.95, time.Second*57),
		Entry("receives arguments (now, now, 1[m], 0.95)", now, now, time.Minute, 0.95, time.Second*57),
		Entry("receives arguments (newYear, newYear-45[m], 1[h], 0.95)", newYear, newYear.Add(-45*time.Minute), time.Hour, 0.95, time.Nanosecond*719999999999),
	)
})
