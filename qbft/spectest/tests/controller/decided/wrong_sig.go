package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongSignature tests a single decided received with a wrong signature
func WrongSignature() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return tests.NewControllerSpecTest(
		"decide wrong sig",
		"Test a single decided message received with a wrong signature, expecting validation error.",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[4]}, []types.OperatorID{1, 2, 3}),
				},
			},
		},
		nil,
		"invalid decided msg: invalid decided msg: msg signature invalid: crypto/rsa: verification error",
		nil,
	)
}
