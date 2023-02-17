package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1FutureMsgs tests a f+1 future msgs that trigger decided sync
func F1FutureMsgs() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "f+1 future msgs",
		InputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     5,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     10,
				Round:      3,
				Identifier: identifier[:],
				Root:       testingutils.TestingQBFTRootData,
			}),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "98c0625374dc64e6350eb704812d9222f9a6121a87c2a55a8b1a3f8790e87c77",
	}
}
