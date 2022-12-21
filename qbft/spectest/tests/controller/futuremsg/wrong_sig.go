package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSig tests future msg with invalid sig
func WrongSig() *ControllerSyncSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "future msg wrong sig",
		InputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[3], 2, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     10,
				Round:      3,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "7b74be21fcdae2e7ed495882d1a499642c15a7f732f210ee84fb40cc97d1ce96",
		ExpectedError:        "invalid future msg: commit msg signature invalid: failed to verify signature",
	}
}
