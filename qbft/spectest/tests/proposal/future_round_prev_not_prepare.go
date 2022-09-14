package proposal

// TODO<olegshmuelov>: PROPOSAL fix test
// FutureRoundPrevNotPrepared tests a proposal for future round, currently not prepared
/*func FutureRoundPrevNotPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = qbft.FirstRound

	rcMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, rcMsgs, nil),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal future round prev not prepared",
		Pre:           pre,
		PostRoot:      "b261accce13a88f020c601d0314bd6eaecd6cb0cea3232198b258cdcc55c1263",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      10,
				Identifier: []byte{1, 2, 3, 4},
				Input: []byte{1, 2, 3, 4},
			}),
		},
	}
}*/
