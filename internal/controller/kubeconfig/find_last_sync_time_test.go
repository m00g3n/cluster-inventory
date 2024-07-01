package kubeconfig

import (
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

var _ = Describe("findLastSyncTime", func() {

	DescribeTable("should return expected values when",
		func(annotations map[string]string, expectedFound bool, expectedTime time.Time) {
			lastSyncTime, found := findLastSyncTime(annotations)
			Expect(found).To(Equal(expectedFound))
			Expect(lastSyncTime).To(Equal(expectedTime))
		},
		Entry("receives empty annotation map", make(map[string]string), false, time.Time{}),
		Entry("receives annotation map containing valid date value",
			map[string]string{lastKubeconfigSyncAnnotation: "2023-01-01T12:00:00Z"}, true,
			func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
				return t
			}()),
		Entry("receives annotation map containing invalid date value", map[string]string{lastKubeconfigSyncAnnotation: "invalid"}, false, time.Time{}),
	)
})
