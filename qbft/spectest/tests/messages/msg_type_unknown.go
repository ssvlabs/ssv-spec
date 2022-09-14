package messages

// MsgTypeUnknown TODO<olegshmuelov> validate message type for unknown or non-exist
// MsgTypeUnknown tests Message type > 5
//func MsgTypeUnknown() *tests.MsgSpecTest {
//	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		MsgType:    6,
//		Height:     qbft.FirstHeight,
//		Round:      qbft.FirstRound,
//		Identifier: []byte{1, 2, 3, 4},
//		Input: []byte{1, 2, 3, 4},
//	})
//
//	return &tests.MsgSpecTest{
//		Name: "msg type unknown",
//		Messages: []*qbft.SignedMessage{
//			msg,
//		},
//		ExpectedError: "message type is invalid",
//	}
//}
