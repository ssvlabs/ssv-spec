package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoPrepareQuorum tests a commit msg received without a previous prepare quorum
func NoPrepareQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		// only 2 out of 4 prepare msgs

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "commit no prepare quorum",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		},
	}
}
