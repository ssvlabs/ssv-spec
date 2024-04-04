package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRound tests prepare for signedProposal.Message.Round > state.Round
func FutureRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare future round",
		Pre:           pre,
		PostRoot:      "e3194c84f99e73171890f32848497b619050587254bf2315ed757095ced37839",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
