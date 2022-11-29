package spectest

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost/blame"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost/keygen"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost/resharing"
)

type SpecTest interface {
	TestName() string
	Run(t *testing.T)
}

var AllTests = []SpecTest{
	tests.HappyFlow(),
	// tests.ResharingHappyFlow(),

	keygen.HappyFlow(),
	resharing.HappyFlow(),
	blame.BlameTypeInvalidCommitment_HappyFlow(),
	blame.BlameTypeInvalidScalar_HappyFlow(),
	blame.BlameTypeInconsistentMessage_HappyFlow(),
	blame.BlameTypeInvalidShare_HappyFlow(),
	blame.BlameTypeInvalidShare_FailedDecrypt_HappyFlow(),
}
