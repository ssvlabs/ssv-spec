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
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))
	pre.ForceStop()

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return tests.NewMsgProcessingSpecTest(
		"force stop round change message",
		"Test round change message processing after instance is force stopped, expecting error.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"instance stopped processing messages",
		nil,
	)
}
