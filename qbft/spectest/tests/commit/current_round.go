package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CurrentRound tests a commit msg with current round, should process
func CurrentRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return tests.NewMsgProcessingSpecTest(
		"commit current round",
		"Test processing of a commit message in the current round",
		pre,
		"",
		nil,
		msgs,
		nil,
		"",
		nil,
	)
}
