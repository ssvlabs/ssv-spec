package decided

// ImparsableData tests a decided msg received with the wrong commit data
// TODO<olegshmuelov>: not relevant anymore
/*func ImparsableData() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide imparsable data",
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
							Data:       []byte{1, 2, 3, 4},
						}),
				},
				ControllerPostRoot: "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
			},
		},
		ExpectedError: "invalid decided msg: invalid decided msg: could not get msg commit data: could not decode commit data from message: invalid character '\\x01' looking for beginning of value",
	}
}*/
