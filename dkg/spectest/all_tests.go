package spectest

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),

	frost.Keygen(),
	frost.Resharing(),
	frost.BlameTypeInvalidShare(),
	frost.BlameTypeInconsistentMessage(),
	tests.ResharingHappyFlow(),
}
