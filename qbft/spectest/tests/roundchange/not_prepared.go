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

	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  nil,
	})
	rcMsgEncoded, _ := rcMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "round change not prepared",
		Pre:              pre,
		PostRoot:         "55cf35ed339dc8b6ee2dbd4ae3af7509dc6305d64252d3d3167fe28a860a6f32",
		InputMessagesSIP: msgs,
		OutputMessages:   []*qbft.SignedMessage{},
	}
}
