package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotProposer tests a justified round change but node is not the proposer
func NotProposer() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Height = tests.ChangeProposerFuncInstanceHeight // will change proposer default for tests

	rcMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: tests.ChangeProposerFuncInstanceHeight,
		Round:  2,
		Input:  &qbft.Data{},
	}).Encode()
	rcMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: tests.ChangeProposerFuncInstanceHeight,
		Round:  2,
		Input:  &qbft.Data{},
	}).Encode()
	rcMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: tests.ChangeProposerFuncInstanceHeight,
		Round:  2,
		Input:  &qbft.Data{},
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded3,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change justification not proposer",
		Pre:           pre,
		PostRoot:      "25675f6c1ffdb51a4539bce513ce71f25093d929971e1a5f5da6a455b8592537",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
				Data: rcMsgEncoded,
			},
		},
	}
}
