package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ThirteenOperators tests a simple full happy flow until decided
func ThirteenOperators() *MsgProcessingSpecTest {
	pre := testingutils.ThirteenOperatorsInstance()
	baseMsgId := types.NewBaseMsgID(testingutils.Testing13SharesSet().ValidatorPK.Serialize(), types.BNRoleAttester)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded4, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded5, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[5], types.OperatorID(5), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded6, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[6], types.OperatorID(6), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded7, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[7], types.OperatorID(7), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded8, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[8], types.OperatorID(8), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded9, _ := testingutils.SignQBFTMsg(testingutils.Testing13SharesSet().Shares[9], types.OperatorID(9), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusProposeMsgType),
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded3,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded4,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded5,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded6,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded7,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded8,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded9,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded3,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded4,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded5,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded6,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded7,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded8,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded9,
		},
	}

	return &MsgProcessingSpecTest{
		Name:          "happy flow thirteen operators",
		Pre:           pre,
		PostRoot:      "19ca526d736aa0f589085724e838d8f6c36c4bc562065d70fdebdfde0a769e54",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(baseMsgId, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
				Data: signMsgEncoded,
			},
		},
	}
}
