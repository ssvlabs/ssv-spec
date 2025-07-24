package tests

import (
	testdoc "github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TenOperators tests a simple full happy flow until decided
func TenOperators() SpecTest {
	pre := testingutils.TenOperatorsInstance()
	ks := testingutils.Testing10SharesSet()

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[4], types.OperatorID(4)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[5], types.OperatorID(5)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[6], types.OperatorID(6)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[7], types.OperatorID(7)),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[4], types.OperatorID(4)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[5], types.OperatorID(5)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[6], types.OperatorID(6)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[7], types.OperatorID(7)),
	}

	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	test := NewMsgProcessingSpecTest(
		"happy flow ten operators",
		testdoc.MsgProcessingHappyFlowTenOperatorsDoc,
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		"",
		nil,
		ks,
	)

	return test
}
