package futuremsg

// WrongSig tests future msg with invalid sig
/*func WrongSig() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "future msg wrong sig",
		InputMessages: []*types.Message{
			testingutils.SignQBFTMsg(ks.Shares[3], 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     10,
				Round:      3,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
		ExpectedError:        "invalid future msg: commit msg signature invalid: failed to verify signature",
	}
}*/
