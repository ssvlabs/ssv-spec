package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PrepareQuorumTriggeredTwiceLateCommit tests triggering prepare quorum twice by sending > 2f+1 prepare messages.
// The commit message is processed after the second prepare quorum is triggered.
func PrepareQuorumTriggeredTwiceLateCommit() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	sc := prepareQuorumTriggeredTwiceStateComparison()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[4], 4),
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepared quorum committed twice late commit",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		},
	}
}
