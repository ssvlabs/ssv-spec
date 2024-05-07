package tests

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TenOperators tests a simple full happy flow until decided
func TenOperators() SpecTest {
	pre := testingutils.TenOperatorsInstance()
	ks := testingutils.Testing10SharesSet()

	msgs := []*types.SignedSSVMessage{
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
	return &MsgProcessingSpecTest{
		Name:          "happy flow ten operators",
		Pre:           pre,
		PostRoot:      "b254e9886acfc2fd4cbb4f08d0d59821a767cdc284eda32c107a5ee7a78358cc",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		},
	}
}
