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

	return &ControllerSyncSpecTest{
		Name: "future msgs multiple signers",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingCommitMultiSignerMessageWithHeight(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				2,
			),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
		ExpectedError:        "invalid future msg: allows 1 signer",
	}
}
