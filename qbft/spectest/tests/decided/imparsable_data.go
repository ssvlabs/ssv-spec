package decided

// TODO<olegshmuelov>: irrelevant test
// ImparsableData tests a decided msg received with the wrong commit data
/*func ImparsableData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Input: []byte{1, 2, 3, 4}
		})
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

		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       []byte{1, 2, 3, 4},
			}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "decided imparsable data",
		Pre:           pre,
		PostRoot:      "a272dbf34be030245fcc44b3210f3137e0cc47e745d0130584de7ff17a47123f",
		InputMessages: msgs,
		ExpectedError: "invalid decided msg: invalid decided msg: could not get msg commit data: could not decode commit data from message: invalid character '\\x01' looking for beginning of value",
		OutputMessages: []*types.Message{},
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
