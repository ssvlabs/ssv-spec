package spectest

import (
	"github.com/bloxapp/ssv-spec/gg20/spectest/tests"
	"testing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),
	tests.HappyFlowNonContinuous(),
	tests.SevenOperators(),
	tests.TenOperators(),
	tests.ThirteenOperators(),
}
