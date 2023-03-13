package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSig tests future msg with invalid sig
func WrongSig() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future msg wrong sig",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 2, 3, 10,
				identifier[:], testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "3b9cd21ca426a4e9e3188e0c8d931861a8f263636c4c0369da84fe9a99fb2fa5",
		ExpectedError:        "invalid future msg: msg signature invalid: failed to verify signature",
	}
}
