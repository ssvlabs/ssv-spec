package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMsg tests future msg invalid msg
func InvalidMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	msg := testingutils.TestingPrepareMessageWithParams(
		ks.Shares[3], 3, 3, 10, identifier[:], testingutils.TestingQBFTRootData)
	msg.Signers = nil

	return &ControllerSyncSpecTest{
		Name: "future invalid msg",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "3b9cd21ca426a4e9e3188e0c8d931861a8f263636c4c0369da84fe9a99fb2fa5",
		ExpectedError:        "invalid future msg: invalid decided msg: message signers is empty",
	}
}
