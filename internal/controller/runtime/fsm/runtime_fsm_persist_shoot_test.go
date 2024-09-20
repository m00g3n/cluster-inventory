package fsm

import (
	"bytes"
	"context"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"io"
	"time"

	"github.com/kyma-project/infrastructure-manager/internal/controller/runtime/fsm/testing"
	"github.com/kyma-project/infrastructure-manager/internal/gardener/shoot"
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

	expectedRuntime := testing.RuntimeOnlyName.DeepCopy()
	expectedRuntime.Spec.Shoot.Provider.Type = "aws"

	It("should persist shoot data", func() {
		next, _, err := sFnDumpShootSpec(testCtx, must(newFakeFSM, withStorageWriter(testWriterGetter), withConverterConfig(shoot.ConverterConfig{})), &systemState{shoot: &testing.ShootNoDNS, instance: *expectedRuntime})
		Expect(err).To(BeNil())
		Expect(next).To(haveName("sFnUpdateStatus"))

		var shootStored gardener.Shoot

		err = yaml.Unmarshal(b.Bytes(), &shootStored)
		Expect(err).To(BeNil())
		Expect(shootStored.ObjectMeta.CreationTimestamp).To(Not(Equal(time.Time{})))
	})
})
