package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// DuplicateSigners tests a multi signer commit msg with duplicate signers
func DuplicateSigners() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)
	commit := testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2]}, []types.OperatorID{1, 2})
	commit.Signers = []types.OperatorID{1, 1}

	return &tests.MsgProcessingSpecTest{
		Name:     "duplicate signers",
		Pre:      pre,
		PostRoot: "470d1a88e97b20eafb08ad9682c10642de27515fff7a8ef3c2d2e97953432357",
		InputMessages: []*qbft.SignedMessage{
			commit,
		},
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: invalid signed message: non unique signer",
	}
}
