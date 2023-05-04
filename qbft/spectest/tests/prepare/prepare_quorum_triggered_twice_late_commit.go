package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PrepareQuorumTriggeredTwiceLateCommit is a spec test that checks the case where all prepares are sent and quorum event is triggered more than once.
// A commit message was seen only after the last prepare
func PrepareQuorumTriggeredTwiceLateCommit() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	sc := prepareQuorumTriggeredTwiceStateComparison()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingPrepareMessage(ks.Shares[4], 4),
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepared quorum committed twice late commit",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
			// ISSUE 214: we should have only commit broadcasted
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
	}
}
