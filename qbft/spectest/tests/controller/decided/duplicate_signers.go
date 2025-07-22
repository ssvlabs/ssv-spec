package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateSigners tests a decided msg with duplicate signers
func DuplicateSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 10)
	msg.OperatorIDs = []types.OperatorID{1, 2, 2}

	return tests.NewControllerSpecTest(
		"decide duplicate signer",
		testdoc.ControllerDecidedDuplicateSignersDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					msg,
				},
			},
		},
		nil,
		"invalid decided msg: invalid decided msg: signed commit invalid: invalid SignedSSVMessage: non unique signer",
		nil,
	)
}
