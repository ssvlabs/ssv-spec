package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1FutureMsgs tests a f+1 future msgs that trigger decdied futuremsg
func F1FutureMsgs() *ControllerSyncSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "f+1 future msgs",
		InputMessages: []*qbft.SignedMessage{
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
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "525bd132a2301e0895b2ad752a34075478f6072dd1fec8ff526833f6342f4ba7",
	}
}
