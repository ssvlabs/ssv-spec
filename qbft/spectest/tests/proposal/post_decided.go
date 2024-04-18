package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		PostRoot:      "0435531b93a3c21862b48c2d18fb85e7ab7539f738fb15ced288ce15b6193974",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks10.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks10.Shares[1], types.OperatorID(1)),
		},
		ExpectedError: "invalid signed message: proposal is not valid with current state",
	}
}
