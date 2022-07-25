package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotPrepared tests a round change msg for non-prepared state
func NotPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound, []byte{1, 2, 3, 4}),
		}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change not prepared",
		Pre:            pre,
		PostRoot:       "26b54c2cc1b6f0cf3b7e49345afc6dc2211f3b895831ed7d383de38650b364e5",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
