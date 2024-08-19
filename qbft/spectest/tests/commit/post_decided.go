package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostDecided tests processing a commit msg after instance decided
func PostDecided() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], 3),
		testingutils.TestingCommitMessage(ks.OperatorKeys[4], 4),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "post decided",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
