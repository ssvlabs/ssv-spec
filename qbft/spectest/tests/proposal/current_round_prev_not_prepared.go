package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CurrentRoundPrevNotPrepared tests a > first round proposal prev not prepared
func CurrentRoundPrevNotPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 10),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(
			ks.Shares[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal happy flow round > 1 (prev not prepared)",
		Pre:           pre,
		PostRoot:      "5901c834260f6533c0c480ee669d45a2d552dcf289c64d0a00f8ebfc38d4f604",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 10),
		},
	}
}
