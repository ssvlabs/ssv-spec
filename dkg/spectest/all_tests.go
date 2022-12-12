package spectest

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),
	// tests.ResharingHappyFlow(),
}
