package tests

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// TenOperators tests a simple full happy flow until decided
func TenOperators() SpecTest {
	pre := testingutils.TenOperatorsInstance()
	ks := testingutils.Testing10SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[4], types.OperatorID(4)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[5], types.OperatorID(5)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[6], types.OperatorID(6)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[7], types.OperatorID(7)),

		testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[4], types.OperatorID(4)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[5], types.OperatorID(5)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[6], types.OperatorID(6)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[7], types.OperatorID(7)),
	}
	return &MsgProcessingSpecTest{
		Name:          "happy flow ten operators",
		Pre:           pre,
		PostRoot:      "5343638957c46929c846d4f10aa288d039ba22e862e4a7a8ed744253819387d6",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		},
	}
}
