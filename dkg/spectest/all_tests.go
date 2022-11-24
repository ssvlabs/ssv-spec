package spectest

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2/blame"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2/keygen"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2/resharing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),
	// tests.ResharingHappyFlow(),

	// frost.Keygen(),
	// frost.Resharing(),
	// frost.BlameTypeInvalidCommitment(),
	// frost.BlameTypeInvalidScalar(),
	// frost.BlameTypeInvalidShare_FailedShareDecryption(),
	// frost.BlameTypeInvalidShare_FailedValidationAgainstCommitment(),
	// frost.BlameTypeInconsistentMessage(),

	keygen.HappyFlow(),
	resharing.HappyFlow(),
	blame.BlameTypeInvalidCommitment_HappyFlow(),
	blame.BlameTypeInvalidScaler_HappyFlow(),
	blame.BlameTypeInconsistentMessage_HappyFlow(),
	blame.BlameTypeInvalidShare_HappyFlow(),
	blame.BlameTypeInvalidShare_FailedDecrypt_HappyFlow(),
}
