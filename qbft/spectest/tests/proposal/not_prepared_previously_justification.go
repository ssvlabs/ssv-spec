package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotPreparedPreviouslyJustification tests a proposal for > 1 round, not prepared previously
func NotPreparedPreviouslyJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 5
	ks := testingutils.Testing4SharesSet()
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 5),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 5,
			testingutils.MarshalJustifications(rcMsgs)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal justification (not prepared)",
		Pre:           pre,
		PostRoot:      "339b6b43d39cd8163e56b9721edfc66a0723c1e8921343571d65f80a4668b574",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5),
		},
	}
}
