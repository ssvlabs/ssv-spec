package roundchange

// F1DuplicateSigner tests not accepting f+1 speed for duplicate signer
//func F1DuplicateSigner() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//
//	prepareMsgs := []*qbft.SignedMessage{
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.PrepareMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Input: []byte{1, 2, 3, 4},
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//			MsgType:    qbft.PrepareMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Input: []byte{1, 2, 3, 4},
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//			MsgType:    qbft.PrepareMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      qbft.FirstRound,
//			Identifier: []byte{1, 2, 3, 4},
//			Input: []byte{1, 2, 3, 4},
//		}),
//	}
//
//	msgs := []*qbft.SignedMessage{
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      2,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      10,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
//		}),
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:           "round change f+1 duplicate",
//		Pre:            pre,
//		PostRoot:       "cc38402bbf897a098b8c96c0391b2c0053bf2663b143d1529151d607b92b610e",
//		InputMessages:  msgs,
//		OutputMessages: []*qbft.SignedMessage{},
//	}
//}

//func FutureRound() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//
//	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	})
//	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	})
//	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	})
//	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  2,
//		Input:  nil,
//	})
//	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  5,
//		Input:  nil,
//	})
//	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height:        qbft.FirstHeight,
//		Round:         10,
//		Input:         []byte{1, 2, 3, 4},
//		PreparedRound: qbft.FirstRound,
//	})
//
//	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessageHeader()
//	prepareMsgHeader2, _ := signQBFTMsg2.ToSignedMessageHeader()
//	prepareMsgHeader3, _ := signQBFTMsg3.ToSignedMessageHeader()
//
//	prepareJustifications := []*qbft.SignedMessageHeader{
//		prepareMsgHeader,
//		prepareMsgHeader2,
//		prepareMsgHeader3,
//	}
//	rcMsg3.RoundChangeJustifications = prepareJustifications
//
//	rcMsgEncoded, _ := rcMsg.Encode()
//	rcMsgEncoded2, _ := rcMsg2.Encode()
//	rcMsgEncoded3, _ := rcMsg3.Encode()
//
//	msgs := []*types.Message{
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: rcMsgEncoded,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: rcMsgEncoded2,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: rcMsgEncoded3,
//		},
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:             "round change future round",
//		Pre:              pre,
//		PostRoot:         "8fb6539597b7fd80818b641fb831e9d3fe8258a44efe0095ec212817e447e1ff",
//		InputMessagesSIP: msgs,
//		OutputMessages:   []*qbft.SignedMessage{},
//	}
//}
