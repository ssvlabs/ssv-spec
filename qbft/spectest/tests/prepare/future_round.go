package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FutureRound tests prepare for signedProposal.Message.Round > state.Round
func FutureRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 3),
	}

	return tests.NewMsgProcessingSpecTest(
		"prepare future round",
		testdoc.PrepareFutureRoundDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: wrong msg round",
		nil,
	)
}
