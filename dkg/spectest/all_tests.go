package spectest

import (
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),
}
