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
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change not prepared",
		Pre:            pre,
		PostRoot:       "d505d40779a07ddd392d4da88a33a2bf67775777eb5227a53cdd3421d5efa28a",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
