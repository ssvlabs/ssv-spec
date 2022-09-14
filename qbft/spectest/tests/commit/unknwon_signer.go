package commit

// UnknownSigner tests a single commit received with an unknown signer
/*func UnknownSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Input: []byte{1, 2, 3, 4},
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Input: []byte{1, 2, 3, 4},
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Input: []byte{1, 2, 3, 4},
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(5), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Input: []byte{1, 2, 3, 4},
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "unknown commit signer",
		Pre:           pre,
		PostRoot:      "a272dbf34be030245fcc44b3210f3137e0cc47e745d0130584de7ff17a47123f",
		InputMessages: msgs,
		ExpectedError: "commit msg invalid: invalid commit msg: commit msg signature invalid: unknown signer",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Input: []byte{1, 2, 3, 4},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Input: []byte{1, 2, 3, 4},
			}),
		},
	}
}*/

//func FutureDecided() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	}).Encode()
//	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	}).Encode()
//	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	}).Encode()
//	multiSignMsgEncoded, _ := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]}, []types.OperatorID{1, 2, 3}, &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	}).Encode()
//	msgs := []*types.Message{
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
//			Data: signMsgEncoded,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
//			Data: signMsgEncoded,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
//			Data: signMsgEncoded2,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
//			Data: signMsgEncoded3,
//		},
//		{
//			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
//			Data: multiSignMsgEncoded,
//		},
//	}
//	return &tests.MsgProcessingSpecTest{
//		Name:             "future decided",
//		Pre:              pre,
//		PostRoot:         "5eb03a8f4d053b3e3f37b88b13c32aaa86e02253731734580b8fa956dec9c53a",
//		InputMessagesSIP: msgs,
//		OutputMessagesSIP: []*types.Message{
//			{
//				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
//				Data: signMsgEncoded,
//			},
//			{
//				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
//				Data: signMsgEncoded,
//			},
//		},
//	}
//}
