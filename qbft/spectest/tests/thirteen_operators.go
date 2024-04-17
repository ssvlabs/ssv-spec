package tests

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ThirteenOperators tests a simple full happy flow until decided
func ThirteenOperators() SpecTest {
	pre := testingutils.ThirteenOperatorsInstance()
	ks := testingutils.Testing13SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.Shares[4], types.OperatorID(4)),
		testingutils.TestingPrepareMessage(ks.Shares[5], types.OperatorID(5)),
		testingutils.TestingPrepareMessage(ks.Shares[6], types.OperatorID(6)),
		testingutils.TestingPrepareMessage(ks.Shares[7], types.OperatorID(7)),
		testingutils.TestingPrepareMessage(ks.Shares[8], types.OperatorID(8)),
		testingutils.TestingPrepareMessage(ks.Shares[9], types.OperatorID(9)),

		testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.Shares[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.Shares[4], types.OperatorID(4)),
		testingutils.TestingCommitMessage(ks.Shares[5], types.OperatorID(5)),
		testingutils.TestingCommitMessage(ks.Shares[6], types.OperatorID(6)),
		testingutils.TestingCommitMessage(ks.Shares[7], types.OperatorID(7)),
		testingutils.TestingCommitMessage(ks.Shares[8], types.OperatorID(8)),
		testingutils.TestingCommitMessage(ks.Shares[9], types.OperatorID(9)),
	}
	return &MsgProcessingSpecTest{
		Name:          "happy flow thirteen operators",
		Pre:           pre,
		PostRoot:      "9479293587bb018ef2e8ed4eb9b86ed1b7a1a1c1f3cb017a41c4896d87f07082",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}
}
