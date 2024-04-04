package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ForceStop tests processing a round change msg when instance force stopped
func ForceStop() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))
	pre.ForceStop()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessage(ks.Shares[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "force stop round change message",
		Pre:           pre,
		PostRoot:      "e3194c84f99e73171890f32848497b619050587254bf2315ed757095ced37839",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
