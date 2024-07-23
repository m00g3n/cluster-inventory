package shoot_test

import (
	"testing"

	"github.com/onsi/ginkgo/v2" //nolint:revive
	"github.com/onsi/gomega"    //nolint:revive
)

func TestMatcher(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "shoot matcher")
}
