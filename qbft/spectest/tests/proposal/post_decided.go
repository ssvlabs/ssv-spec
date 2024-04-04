package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing proposal msg after instance decided
func PostDecided() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks4 := testingutils.Testing4SharesSet()
	ks10 := testingutils.Testing10SharesSet() // TODO should be 4?

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks4.Shares[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks4.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks4.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks4.Shares[3], types.OperatorID(3)),

		testingutils.TestingCommitMessage(ks4.Shares[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks4.Shares[2], types.OperatorID(2)),
		testingutils.TestingCommitMessage(ks4.Shares[3], types.OperatorID(3)),

		testingutils.TestingProposalMessage(ks4.Shares[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal post decided",
		Pre:           pre,
		PostRoot:      "4b4536bf7f996c924ee4fcc5d4dc05c3b5249b1970f6e5ecb2d6212df1eaf602",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks10.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks10.Shares[1], types.OperatorID(1)),
		},
		ExpectedError: "invalid signed message: proposal is not valid with current state",
	}
}
