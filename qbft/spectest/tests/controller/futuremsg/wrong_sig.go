package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSig tests future msg with invalid sig
func WrongSig() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "future msg wrong sig",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 2, 3, 10,
				testingutils.DefaultIdentifier, testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
		ExpectedError:        "invalid future msg: msg signature invalid: failed to verify signature",
	}
}
