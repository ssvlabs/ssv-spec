package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ThirteenOperators tests a simple full happy flow until decided
func ThirteenOperators() *MsgProcessingSpecTest {
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
		PostRoot:      "ce494f9e2fc81dab675feaab631e7e828424971aec32d28b6b01977973ee4d93",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}
}
