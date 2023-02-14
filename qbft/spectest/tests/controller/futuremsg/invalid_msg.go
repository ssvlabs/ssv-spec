package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMsg tests future msg invalid msg
func InvalidMsg() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10, testingutils.TestingQBFTRootData)
	msg.Signature = nil

	return &ControllerSyncSpecTest{
		Name: "future invalid msg",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
		ExpectedError:        "invalid future msg: invalid decided msg: message signature is invalid",
	}
}
