package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
		PostRoot:      "fc68d10ca3bc1250da3aa789d4feae5fa197d08964c26c9f16cdfd39931d38ce",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
