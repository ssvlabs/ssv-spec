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
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 5),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.NetworkKeys[1], types.OperatorID(1), 5,
			testingutils.MarshalJustifications(rcMsgs)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal justification (not prepared)",
		Pre:           pre,
		PostRoot:      "c36727f091866a125ef815fe77a50df5e38335a1bb44eca2dbe0d2a4c34de4d8",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 5),
		},
	}
}
