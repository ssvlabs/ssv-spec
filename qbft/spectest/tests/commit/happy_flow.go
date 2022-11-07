package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a quorum of commits received after prepare quorum
func HappyFlow() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue).Encode()
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()

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
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded3,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit happy flow",
		Pre:           pre,
		PostRoot:      "388136ceb23b210e0e0e0004d1e7030734e4581ab1a29358bab6031ba7ecffa4",
		InputMessages: msgs,
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
