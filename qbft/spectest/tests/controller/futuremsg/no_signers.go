package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSigners tests future msg with no signers
func NoSigners() *ControllerSyncSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     10,
		Round:      3,
		Identifier: identifier[:],
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})
	msg.Signers = []types.OperatorID{}

	return &ControllerSyncSpecTest{
		Name: "future msgs no signer",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "7b74be21fcdae2e7ed495882d1a499642c15a7f732f210ee84fb40cc97d1ce96",
		ExpectedError:        "invalid future msg: invalid decided msg: message signers is empty",
	}
}
