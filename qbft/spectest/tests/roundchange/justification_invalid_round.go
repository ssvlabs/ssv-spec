package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// JustificationInvalidRound tests justifications with prepared round > msg round
func JustificationInvalidRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 3),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 3),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "justification invalid round",
		Pre:            pre,
		PostRoot:       "96e6d7bdbb98a2d9937f3d97d6aa096bd3a58f923b61012048ac70ad52765919",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: round change justification invalid: wrong msg round",
	}
}
