package messages

// MsgNilIdentifier TODO<olegshmuelov> find a way to validate the identifier
// MsgNilIdentifier tests Message with Identifier == nil
//func MsgNilIdentifier() *tests.MsgSpecTest {
//	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		MsgType:    qbft.CommitMsgType,
//		Height:     qbft.FirstHeight,
//		Round:      qbft.FirstRound,
//		Identifier: nil,
//		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
//	})
//
//	return &tests.MsgSpecTest{
//		Name: "msg identifier nil",
//		Messages: []*qbft.SignedMessage{
//			msg,
//		},
//		ExpectedError: "message identifier is invalid",
//	}
//}
