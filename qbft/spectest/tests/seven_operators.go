package tests

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SevenOperators tests a simple full happy flow until decided
func SevenOperators() SpecTest {
	pre := testingutils.SevenOperatorsInstance()
	ks := testingutils.Testing7SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[4], types.OperatorID(4)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[5], types.OperatorID(5)),

		testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[4], types.OperatorID(4)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[5], types.OperatorID(5)),
	}
	return &MsgProcessingSpecTest{
		Name:          "happy flow seven operators",
		Pre:           pre,
		PostRoot:      "0230e78f218fbc7f8560d0997e9ae003e90d4e1cf710be1b3917cf659ef955c4",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		},
	}
}
