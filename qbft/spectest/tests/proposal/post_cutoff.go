package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostCutoff tests processing a proposal msg when round >= cutoff
func PostCutoff() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 15

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 15),
	}

	return tests.NewMsgProcessingSpecTest(
		"round cutoff proposal message",
		"Test processing proposal message when the round is greater than or equal to the cutoff, expecting the instance to stop processing messages.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"instance stopped processing messages",
		nil,
	)
}
