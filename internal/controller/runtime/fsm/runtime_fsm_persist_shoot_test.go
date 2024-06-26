package fsm

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/kyma-project/infrastructure-manager/internal/controller/runtime/fsm/testing"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
	"sigs.k8s.io/yaml"
)

var _ = Describe("KIM sFnPersist", func() {

	var b bytes.Buffer
	testWriterGetter := func() writerGetter {
		return func(string) (io.Writer, error) {
			return &b, nil
		}
	}()

	testCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	expectedData, err := yaml.Marshal(&testing.ShootNoDNS)
	Expect(err).ShouldNot(HaveOccurred())

	It("shoutld persist shoot data", func() {
		next, _, err := sFnPersistShoot(testCtx, must(newFakeFSM, withStorageWriter(testWriterGetter)), &systemState{shoot: &testing.ShootNoDNS})
		Expect(err).To(BeNil())
		Expect(next).To(haveName("sFnUpdateStatus"))
		Expect(expectedData).To(Equal(b.Bytes()))
	})
})
