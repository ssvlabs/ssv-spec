package spectest

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2/keygen"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	// tests.HappyFlow(),

	frost.Keygen(),
	frost.Resharing(),
	frost.BlameTypeInvalidCommitment(),
	frost.BlameTypeInvalidScalar(),
	frost.BlameTypeInvalidShare_FailedShareDecryption(),
	frost.BlameTypeInvalidShare_FailedValidationAgainstCommitment(),
	frost.BlameTypeInconsistentMessage(),
	tests.ResharingHappyFlow(),

	keygen.HappyFlow(),
}
