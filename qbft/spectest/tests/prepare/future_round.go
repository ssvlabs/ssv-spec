package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FutureRound tests prepare for signedProposal.Message.Round > state.Round
func FutureRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare future round",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
