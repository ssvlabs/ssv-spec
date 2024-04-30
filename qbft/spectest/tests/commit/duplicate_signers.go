package commit

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateSigners tests a multi signer commit msg with duplicate signers
func DuplicateSigners() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1)
	commit := testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]}, []types.OperatorID{1, 2})
	commit.OperatorIDs = []types.OperatorID{1, 1}

	return &tests.MsgProcessingSpecTest{
		Name:     "duplicate signers",
		Pre:      pre,
		PostRoot: "fc68d10ca3bc1250da3aa789d4feae5fa197d08964c26c9f16cdfd39931d38ce",
		InputMessages: []*types.SignedSSVMessage{
			commit,
		},
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: invalid SignedSSVMessage: non unique signer",
	}
}
