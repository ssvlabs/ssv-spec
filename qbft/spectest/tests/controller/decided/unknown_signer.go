package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownSigner tests a decided msg with an unknown signer
func UnknownSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return tests.NewControllerSpecTest(
		"decide unknown signer",
		"Test a decided message with an unknown signer, expecting validation error.",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 5}),
				},
			},
		},
		nil,
		"invalid decided msg: invalid decided msg: signer not in committee",
		nil,
	)
}
