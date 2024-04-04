package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SevenOperators tests a simple full happy flow until decided
func SevenOperators() SpecTest {
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
		PostRoot:      "6cd50bd44cc5a745876ccd2fd9dee748bdcca4a609c719930d9f95c468ed2e65",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}
}
