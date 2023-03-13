package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests future msg signed by unknown signer
func UnknownSigner() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()

	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	msg := testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10,
		identifier[:], testingutils.TestingQBFTRootData)
	msg.Signers = []types.OperatorID{10}

	return &ControllerSyncSpecTest{
		Name: "future msg unknown signer",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "3b9cd21ca426a4e9e3188e0c8d931861a8f263636c4c0369da84fe9a99fb2fa5",
		ExpectedError:        "invalid future msg: msg signature invalid: unknown signer",
	}
}
