package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// DuplicateSigners tests a multi signer commit msg with duplicate signers
func DuplicateSigners() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)
	commit := testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2]}, []types.OperatorID{1, 2})
	commit.Signers = []types.OperatorID{1, 1}

	return &tests.MsgProcessingSpecTest{
		Name:     "duplicate signers",
		Pre:      pre,
		PostRoot: "69c049da1936e3727d09f976754cc7ee3a5cb7d85fa1e079f0465096b0a15cdb",
		InputMessages: []*qbft.SignedMessage{
			commit,
		},
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: invalid signed message: non unique signer",
	}
}
