package decided

// InvalidValCheckData tests a decided message with invalid decided data (but should pass as it's decided)
// TODO<olegshmuelov> what should we check here?
/*func InvalidValCheckData() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide invalid value (should pass)",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     10,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes(testingutils.TestingInvalidValueCheck),
						}),
				},
				SavedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     10,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes(testingutils.TestingInvalidValueCheck),
					}),
				DecidedVal:         testingutils.TestingInvalidValueCheck,
				DecidedCnt:         1,
				ControllerPostRoot: "8be69818570269c47665bcfb8d4a3834387f4ee5d41eaaf03702af5334dd17de",
			},
		},
	}
}*/
