package fsm

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
)

func TestFSM(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "FSM Suite")
}
