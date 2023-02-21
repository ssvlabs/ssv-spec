package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRCJustification tests a proposal for > 1 round, not prepared previously but without quorum of round change msgs justification
func NoRCJustification() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	ks := testingutils.Testing4SharesSet()
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "no rc quorum",
		Pre:            pre,
		PostRoot:       "e9a7c1b25a014f281d084637bdc50ec144f1262b5cc87eeb6b4493d27d10a69b",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: change round has no quorum",
	}
}
