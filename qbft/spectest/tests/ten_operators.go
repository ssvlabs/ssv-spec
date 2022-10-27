package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// TenOperators tests a simple full happy flow until decided
func TenOperators() *MsgProcessingSpecTest {
	pre := testingutils.TenOperatorsInstance()
	baseMsgId := types.NewBaseMsgID(testingutils.Testing10SharesSet().ValidatorPK.Serialize(), types.BNRoleAttester)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded4, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded5, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[5], types.OperatorID(5), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded6, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[6], types.OperatorID(6), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	signMsgEncoded7, _ := testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[7], types.OperatorID(7), &qbft.Message{
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
	}

	return &MsgProcessingSpecTest{
		Name:          "happy flow ten operators",
		Pre:           pre,
		PostRoot:      "6a61afbda62663ca50e75e1d183a92f6e227da8736abbef7c325e7daad828456",
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
