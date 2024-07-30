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

	shootWrite, err := yaml.Marshal(&testing.ShootNoDNS)
	runtimeWrite, err := yaml.Marshal(&testing.RuntimeOnlyName)
	expectedData := append(shootWrite, runtimeWrite...)
	Expect(err).ShouldNot(HaveOccurred())

	It("should persist shoot data", func() {
		next, _, err := sFnDumpShootSpec(testCtx, must(newFakeFSM, withStorageWriter(testWriterGetter)), &systemState{shoot: &testing.ShootNoDNS, instance: testing.RuntimeOnlyName})
		Expect(err).To(BeNil())
		Expect(next).To(haveName("sFnUpdateStatus"))
		Expect(expectedData).To(Equal(b.Bytes()))
	})
})
