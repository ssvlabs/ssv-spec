package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigners tests future msg with multiple signers
func MultiSigners() *ControllerSyncSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future msgs multiple signers",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMultiSignerMessageWithHeightAndIdentifier(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				2,
				identifier[:],
			),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "3b9cd21ca426a4e9e3188e0c8d931861a8f263636c4c0369da84fe9a99fb2fa5",
		ExpectedError:        "invalid future msg: allows 1 signer",
	}
}
