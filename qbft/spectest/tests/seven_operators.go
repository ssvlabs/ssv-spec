package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SevenOperators tests a simple full happy flow until decided
func SevenOperators() *MsgProcessingSpecTest {
	pre := testingutils.SevenOperatorsInstance()
	baseMsgId := types.NewBaseMsgID(testingutils.Testing7SharesSet().ValidatorPK.Serialize(), types.BNRoleAttester)
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue).Encode()
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing7SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing7SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded4, _ := testingutils.SignQBFTMsg(testingutils.Testing7SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded5, _ := testingutils.SignQBFTMsg(testingutils.Testing7SharesSet().Shares[5], types.OperatorID(5), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
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
	}

	return &MsgProcessingSpecTest{
		Name:          "happy flow seven operators",
		Pre:           pre,
		PostRoot:      "31b31a60481c849f72af416d194341e0a1e3c0c810f260fad54aefffac2bc24d",
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
