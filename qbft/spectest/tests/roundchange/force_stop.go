package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ForceStop tests processing a round change msg when instance force stopped
func ForceStop() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))
	pre.ForceStop()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "force stop round change message",
		Pre:           pre,
		PostRoot:      "2253eea5735c33797cd1f1a1e3ced2cb8b16ee1c78ae1747e18041b67216d622",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
