package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostDecided tests processing proposal msg after instance decided
func PostDecided() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks.OperatorKeys[3], types.OperatorID(3)),

		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal post decided",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		},
		ExpectedError: "invalid signed message: proposal is not valid with current state",
	}
}
