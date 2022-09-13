package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsgQuorumPreparedRCFirst tests a duplicate rc msg (the prev prepared one first)
func DuplicateMsgQuorumPreparedRCFirst() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	rcMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change duplicate msg quorum (prev prepared rc first)",
		Pre:           pre,
		PostRoot:      "578e59d3fba60d1cca6e2f022c2845148c2c98dee22aa2a30aa3d6ecce47d768",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      2,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, rcMsgs, prepareMsgs),
			}),
		},
	}
}

//func DuplicateMsgPrepared() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	pre.State.Round = 2
//
//	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	})
//	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	})
//	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	})
//	changeRoundMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  2,
//		Input:  nil,
//	})
//	changeRoundMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height:        qbft.FirstHeight,
//		Round:         2,
//		Input:         []byte{1, 2, 3, 4},
//		PreparedRound: qbft.FirstRound,
//	})
//	changeRoundMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  2,
//		Input:  nil,
//	})
//	changeRoundMsg4 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  2,
//		Input:  nil,
//	})
//	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  2,
//		Input:  []byte{1, 2, 3, 4},
//	})
//
//	prepareMsgHeader, _ := prepareMsg.ToSignedMessageHeader()
//	prepareMsgHeader2, _ := prepareMsg2.ToSignedMessageHeader()
//	prepareMsgHeader3, _ := prepareMsg3.ToSignedMessageHeader()
//
//	changeRoundMsgHeader, _ := changeRoundMsg.ToSignedMessageHeader()
//	changeRoundMsgHeader2, _ := changeRoundMsg2.ToSignedMessageHeader()
//	changeRoundMsgHeader3, _ := changeRoundMsg3.ToSignedMessageHeader()
//	changeRoundMsgHeader4, _ := changeRoundMsg4.ToSignedMessageHeader()
//
//	prepareJustifications := []*qbft.SignedMessageHeader{
//		prepareMsgHeader,
//		prepareMsgHeader2,
//		prepareMsgHeader3,
//	}
//	changeRoundMsg2.RoundChangeJustifications = prepareJustifications
//
//	proposalMsg.RoundChangeJustifications = []*qbft.SignedMessageHeader{
//		changeRoundMsgHeader,
//		changeRoundMsgHeader2,
//		changeRoundMsgHeader3,
//		changeRoundMsgHeader4,
//	}
//	proposalMsg.ProposalJustifications = prepareJustifications
//
//	changeRoundMsgEncoded, _ := changeRoundMsg.Encode()
//	changeRoundMsgEncoded2, _ := changeRoundMsg2.Encode()
//	changeRoundMsgEncoded3, _ := changeRoundMsg3.Encode()
//	changeRoundMsgEncoded4, _ := changeRoundMsg4.Encode()
//	proposalMsgEncoded, _ := proposalMsg.Encode()
//	prepareMsg.Message.Round = 2
//	prepareMsgEncoded, _ := prepareMsg.Encode()
//
//	msgs := []*types.Message{
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: changeRoundMsgEncoded,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: changeRoundMsgEncoded2,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: changeRoundMsgEncoded3,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: changeRoundMsgEncoded4,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
//			Data: proposalMsgEncoded,
//		},
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:             "round change duplicate prepared msg",
//		Pre:              pre,
//		PostRoot:         "9ec4fc8f40d1466da583eef26a5e51af98ee801c5eb6043da544db3d3d523ea0",
//		InputMessagesSIP: msgs,
//		OutputMessagesSIP: []*types.Message{
//			{
//				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
//				Data: proposalMsgEncoded,
//			},
//			{
//				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
//				Data: prepareMsgEncoded,
//			},
//		},
//	}
//}