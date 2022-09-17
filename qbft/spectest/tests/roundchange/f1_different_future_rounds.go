package roundchange

// F1DifferentFutureRounds tests f+1 speedup with one rc prev prepared
//func F1DifferentFutureRounds() *tests.MsgProcessingSpecTest {
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
//			Round:      5,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      10,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
//		}),
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:          "round change f+1 prepared",
//		Pre:           pre,
//		PostRoot:      "c7fbe6d05dd956a638b5a8d6dcaefe5866916bb77d1817e30cbea6d4b3baa172",
//		InputMessages: msgs,
//		OutputMessages: []*qbft.SignedMessage{
//			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//				MsgType:    qbft.RoundChangeMsgType,
//				Height:     qbft.FirstHeight,
//				Round:      5,
//				Identifier: []byte{1, 2, 3, 4},
//				Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
//			}),
//		},
//	}
//}

//func Prepared() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	pre.State.Round = 2
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
//		Height:        qbft.FirstHeight,
//		Round:         2,
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
//	rcMsg.RoundChangeJustifications = prepareJustifications
//
//	rcMsgEncoded, _ := rcMsg.Encode()
//
//	msgs := []*types.Message{
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//			Data: rcMsgEncoded,
//		},
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:             "round change prepared",
//		Pre:              pre,
//		PostRoot:         "3ff2f27d56d3f3f1503509276d2ac8c8cb90688d9dadbe8d16a66302394a46a8",
//		InputMessages: msgs,
//		OutputMessages:   []*qbft.SignedMessage{},
//	}
//}
