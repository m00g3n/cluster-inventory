package fsm

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/kyma-project/infrastructure-manager/internal/controller/runtime/fsm/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/yaml"
)

type writerGetter = func(string) (io.Writer, error)

var _ = Describe("KIM sFnPersist", func() {

	var b bytes.Buffer
	testStringWriter := func() writerGetter {
		return func(string) (io.Writer, error) {
			return &b, nil
		}
	}

	getWriter = testStringWriter()

	testCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	expectedData, err := yaml.Marshal(&testing.ShootNoDNS)
	Expect(err).ShouldNot(HaveOccurred())

	It("shoutld persist shoot data", func() {
		next, _, err := sFnPersistShoot(testCtx, must(newFakeFSM), &systemState{shoot: &testing.ShootNoDNS})
		Expect(err).To(BeNil())
		Expect(next).To(haveName("sFnUpdateStatus"))
		Expect(expectedData).To(Equal(b.Bytes()))
	})
})
