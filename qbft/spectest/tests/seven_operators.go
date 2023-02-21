package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SevenOperators tests a simple full happy flow until decided
func SevenOperators() *MsgProcessingSpecTest {
	pre := testingutils.SevenOperatorsInstance()
	ks := testingutils.Testing7SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.Shares[4], types.OperatorID(4)),
		testingutils.TestingPrepareMessage(ks.Shares[5], types.OperatorID(5)),

		testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.Shares[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.Shares[4], types.OperatorID(4)),
		testingutils.TestingCommitMessage(ks.Shares[5], types.OperatorID(5)),
	}
	return &MsgProcessingSpecTest{
		Name:          "happy flow seven operators",
		Pre:           pre,
		PostRoot:      "53b8d8fc9a5ee6fd2cf3973e4fcbdbbe96a854f9b30287ade46ff080d93e5417",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}
}
