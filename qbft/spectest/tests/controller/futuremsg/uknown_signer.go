package futuremsg

// UnknownSigner tests future msg signed by unknown signer
/*func UnknownSigner() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     10,
		Round:      3,
		Identifier: identifier[:],
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})
	msg.Signers = []types.OperatorID{10}

	return &ControllerSyncSpecTest{
		Name: "future msg unknown signer",
		InputMessages: []*types.Message{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
		ExpectedError:        "invalid future msg: commit msg signature invalid: unknown signer",
	}
}*/
