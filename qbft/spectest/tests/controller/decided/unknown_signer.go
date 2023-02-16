package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// UnknownSigner tests a decided msg with an unknown signer
func UnknownSigner() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide unknown signer",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 5}),
				},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "invalid decided msg: invalid decided msg: msg signature invalid: unknown signer",
	}
}
