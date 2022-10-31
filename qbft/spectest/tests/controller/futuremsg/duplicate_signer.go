package futuremsg

// DuplicateSigner tests multiple future msg for the same signer (doesn't trigger futuremsg)
func DuplicateSigner() *ControllerSyncSpecTest {
	return &ControllerSyncSpecTest{}
	/*identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "future msg duplicate signer",
		InputMessages: []*types.Message{
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     5,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
			testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     10,
				Round:      3,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),

			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     6,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.RoundChangeMsgType,
				Height:     2,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     50,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "4143f41114629c9d7e012ac3ef2b29dafbde78992b8604d50e7c43bb96b027ae",
		ExpectedError:        "discarded future msg",
	}*/
}
