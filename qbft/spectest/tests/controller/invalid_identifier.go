package controller

// InvalidIdentifier tests a process msg with the wrong identifier
//func InvalidIdentifier() *tests.ControllerSpecTest {
//	share := testingutils.Testing4SharesSet().Shares[1]
//	msg := &qbft.Message{
//		MsgType:    qbft.ProposalMsgType,
//		Height:     qbft.FirstHeight,
//		Round:      qbft.FirstRound,
//		Identifier: []byte{1, 2, 3, 4},
//		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
//	}
//	return &tests.ControllerSpecTest{
//		Name: "invalid identifier",
//		RunInstanceData: []struct {
//			InputValue    []byte
//			InputMessages []*qbft.SignedMessage
//			Decided       bool
//			DecidedVal    []byte
//			DecidedCnt    uint
//			SavedDecided  *qbft.SignedMessage
//		}{
//			{
//				InputValue: []byte{1, 2, 3, 4},
//				InputMessages: []*qbft.SignedMessage{
//					testingutils.SignQBFTMsg(share, 1, msg),
//				},
//				Decided:    false,
//				DecidedVal: nil,
//			},
//		},
//		ExpectedError: "message doesn't belong to Identifier",
//	}
//}

//func InvalidIdentifier() *tests.ControllerSpecTest {
//	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
//	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		Height: qbft.FirstHeight,
//		Round:  qbft.FirstRound,
//		Input:  []byte{1, 2, 3, 4},
//	}).Encode()
//	return &tests.ControllerSpecTest{
//		Name: "invalid identifier",
//		RunInstanceData: []struct {
//			InputValue    []byte
//			InputMessages []*types.Message
//			Decided       bool
//			DecidedVal    []byte
//			DecidedCnt    uint
//			SavedDecided  *qbft.SignedMessage
//		}{
//			{
//				InputValue: []byte{1, 2, 3, 4},
//				InputMessages: []*types.Message{
//					{
//						ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
//						Data: signMsgEncoded,
//					}},
//				Decided:    false,
//				DecidedVal: nil,
//			},
//		},
//		ExpectedError: "message doesn't belong to Identifier",
//	}
//}
