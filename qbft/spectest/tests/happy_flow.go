package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
//func HappyFlow() *MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	msgs := []*qbft.SignedMessage{
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.ProposalMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
//		}),
//
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.PrepareMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//			MsgType:    qbft.PrepareMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//			MsgType:    qbft.PrepareMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
//		}),
//
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.CommitMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//			MsgType:    qbft.CommitMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//			MsgType:    qbft.CommitMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
//		}),
//	}
//	return &MsgProcessingSpecTest{
//		Name:          "happy flow",
//		Pre:           pre,
//		PostRoot:      "7a305edc0784ac3a70285e9404d403aac1dd9c5cd4f7b70cac3824d026cc9804",
//		InputMessages: msgs,
//		OutputMessages: []*qbft.SignedMessage{
//			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//				MsgType:    qbft.PrepareMsgType,
//				Height:     qbft.FirstHeight,
//				Round:      qbft.FirstRound,
//				Identifier: []byte{1, 2, 3, 4},
//				Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
//			}),
//			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//				MsgType:    qbft.CommitMsgType,
//				Height:     qbft.FirstHeight,
//				Round:      qbft.FirstRound,
//				Identifier: []byte{1, 2, 3, 4},
//				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
//			}),
//		},
//	}
//}

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	baseMsgId := types.NewBaseMsgID(testingutils.Testing4SharesSet().ValidatorPK.Serialize(), types.BNRoleAttester)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
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
	}
	return &MsgProcessingSpecTest{
		Name:             "happy flow",
		Pre:              pre,
		PostRoot:         "2d11238c88223c7a2dcf161ab1ed04818d4aaa861bfadf890dcc5e1fc6aaef45",
		InputMessagesSIP: msgs,
		OutputMessagesSIP: []*types.Message{
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
