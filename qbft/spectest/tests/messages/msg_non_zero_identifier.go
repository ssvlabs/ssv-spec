package messages

// MsgNonZeroIdentifier TODO<olegshmuelov> find a way to validate for identifier
// MsgNonZeroIdentifier tests Message with len(Identifier) == 0
//func MsgNonZeroIdentifier() *tests.MsgSpecTest {
//	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		MsgType:    qbft.CommitMsgType,
//		Height:     qbft.FirstHeight,
//		Round:      qbft.FirstRound,
//		Identifier: []byte{},
//		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
//	})
//
//	return &tests.MsgSpecTest{
//		Name: "msg identifier len == 0",
//		Messages: []*qbft.SignedMessage{
//			msg,
//		},
//		ExpectedError: "message identifier is invalid",
//	}
//}
