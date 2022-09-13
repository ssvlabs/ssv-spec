package roundchange

// UnknownSigner tests a signed round change msg with an unknown signer
//func UnknownSigner() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	pre.State.Round = 2
//
//	msgs := []*qbft.SignedMessage{
//		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(5), &qbft.Message{
//			MsgType:    qbft.RoundChangeMsgType,
//			Height:     qbft.FirstHeight,
//			Round:      2,
//			Identifier: []byte{1, 2, 3, 4},
//			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
//		}),
//	}
//
//	return &tests.MsgProcessingSpecTest{
//		Name:           "round change unknown signer",
//		Pre:            pre,
//		PostRoot:       "4aafcc4aa9e2435579c85aa26e659fe650aefb8becb5738d32dd9286f7ff27c3",
//		InputMessages:  msgs,
//		OutputMessages: []*qbft.SignedMessage{},
//		ExpectedError:  "round change msg invalid: round change msg signature invalid: unknown signer",
//	}
//}

//func NotPrepared() *tests.MsgProcessingSpecTest {
//	pre := testingutils.BaseInstance()
//	pre.State.Round = 2
//
//	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  2,
//		Input:  nil,
//	})
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
//		Name:             "round change not prepared",
//		Pre:              pre,
//		PostRoot:         "55cf35ed339dc8b6ee2dbd4ae3af7509dc6305d64252d3d3167fe28a860a6f32",
//		InputMessagesSIP: msgs,
//		OutputMessages:   []*qbft.SignedMessage{},
//	}
//}
