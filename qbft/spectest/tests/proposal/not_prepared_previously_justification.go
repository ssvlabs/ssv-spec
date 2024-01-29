package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotPreparedPreviouslyJustification tests a proposal for > 1 round, not prepared previously
func NotPreparedPreviouslyJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 5
	ks := testingutils.Testing4SharesSet()
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 5),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 5,
			testingutils.MarshalJustifications(rcMsgs)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal justification (not prepared)",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 5),
		},
	}
}
