package roundchange

// QuorumPrepared tests a round change msg for prepared state
//func QuorumPrepared() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	pre.State.Round = 2
//
//	prepareMsgs := []*qbft.SignedMessage{
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
//	}
//	msgs := []*qbft.SignedMessage{
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      2,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      2,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
//		}),
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      2,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
//		}),
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:          "round change prepared",
//		Pre:           pre,
//		PostRoot:      "693f301963e027b305656d88af9eeb312f70216c49b16661a8ffce3fc6409e70",
//		InputMessages: msgs,
//		OutputMessages: []*qbft.SignedMessage{
//			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//				MsgType:    qbft.ProposalMsgType,
//				Height:     qbft.FirstHeight,
//				Round:      2,
//				Identifier: []byte{1, 2, 3, 4},
//				Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, msgs, prepareMsgs),
//			}),
//		},
//	}
//}

//func F1SpeedupDifferentRounds() *tests.MsgProcessingSpecTest {
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
//	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  5,
//		Input:  nil,
//	})
//	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//		Height:        qbft.FirstHeight,
//		Round:         10,
//		Input:         []byte{1, 2, 3, 4},
//		PreparedRound: qbft.FirstRound,
//	})
//	outputRcMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  5,
//		Input:  nil,
//	}).Encode()
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
//		Name:             "round change speedup different rounds",
//		Pre:              pre,
//		PostRoot:         "cca2c7f305b3c56956818a56fcca8a51fd7a96a7dda1efdc4214b9a5a29acae1",
//		InputMessagesSIP: msgs,
//		OutputMessagesSIP: []*types.Message{
//			{
//				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//				Data: rcMsgEncoded,
//			},
//			{
//				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
//				Data: outputRcMsgEncoded,
//			},
//		},
//	}
//}
