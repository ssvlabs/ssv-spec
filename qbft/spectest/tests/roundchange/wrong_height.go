package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongHeight tests a round change msg with wrong height
func WrongHeight() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: 2,
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
		Name:           "round change invalid height",
		Pre:            pre,
		PostRoot:       "a8b80879ebf2ecee42fddc69b67dd5f6adfd6aa8b7114246aec80ce1bfef513a",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "round change msg invalid: round change Height is wrong",
	}
}
