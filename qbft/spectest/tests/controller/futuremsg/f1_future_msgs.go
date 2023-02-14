package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1FutureMsgs tests a f+1 future msgs that trigger decided sync
func F1FutureMsgs() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "f+1 future msgs",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingCommitMessageWithHeight(ks.Shares[4], 4, 5),
			testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10, testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 1,
		ControllerPostRoot:   "ca7eaf0f0b404b601dc8dd471924794ce32ef6bcb88721098b7b6014001754c1",
	}
}
