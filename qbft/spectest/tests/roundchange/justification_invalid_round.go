package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// JustificationInvalidRound tests justifications with prepared round > msg round
func JustificationInvalidRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 3),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 3),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 3),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "justification invalid round",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: round change justification invalid: wrong msg round",
	}
}
