package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRound tests prepare for signedProposal.Message.Round > state.Round
func FutureRound() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare future round",
		Pre:           pre,
		PostRoot:      "69c049da1936e3727d09f976754cc7ee3a5cb7d85fa1e079f0465096b0a15cdb",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
