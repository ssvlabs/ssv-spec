package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ThirteenOperators tests a simple full happy flow until decided
func ThirteenOperators() *MsgProcessingSpecTest {
	pre := testingutils.ThirteenOperatorsInstance()
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue).Encode()
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded4, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded5, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[5], types.OperatorID(5), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded6, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[6], types.OperatorID(6), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded7, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[7], types.OperatorID(7), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded8, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[8], types.OperatorID(8), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded9, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[9], types.OperatorID(9), &qbft.Message{
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
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded4,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded5,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded6,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded7,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded8,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded9,
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
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded4,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded5,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded6,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded7,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded8,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded9,
		},
	}

	return &MsgProcessingSpecTest{
		Name:          "happy flow thirteen operators",
		Pre:           pre,
		PostRoot:      "3fe5292a9453e7c02568473f55796e5b6dbb62ce9f743e35309fab62323b29a0",
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
