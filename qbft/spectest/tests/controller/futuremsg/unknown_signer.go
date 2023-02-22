package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests future msg signed by unknown signer
func UnknownSigner() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10,
		testingutils.DefaultIdentifier, testingutils.TestingQBFTRootData)
	msg.Signers = []types.OperatorID{10}

	return &ControllerSyncSpecTest{
		Name: "future msg unknown signer",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
		ExpectedError:        "invalid future msg: msg signature invalid: unknown signer",
	}
}
