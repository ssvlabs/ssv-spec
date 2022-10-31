package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongData2 tests a single commit received with a different commit data than the prepared data
func WrongData2() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  pre.StartValue,
	}).Encode()
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	}).Encode()
	signMsgEncodedInvalid, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 3}},
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded3,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncodedInvalid,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit data != prepared data",
		Pre:           pre,
		PostRoot:      "3cc4269ff493d8ce9b2c9a43690843a4758b9ff77c2d6c1b659b3cc1f6fa6748",
		InputMessages: msgs,
		ExpectedError: "commit msg invalid: proposed data different than commit msg data",
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: signMsgEncoded,
			},
		},
	}
}
