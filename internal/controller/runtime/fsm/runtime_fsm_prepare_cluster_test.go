package fsm

import (
	"context"

	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	_ "github.com/onsi/gomega"
)

var _ = Describe("KIM sFnInitialise", func() {
	_ = metav1.NewTime(time.Now())

	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
})
