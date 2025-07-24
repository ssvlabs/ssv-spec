package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongMsgType tests a non commit msg with 2f+1 signers
func WrongMsgType() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	test := tests.NewControllerSpecTest(
		"decide wrong msg type",
		testdoc.ControllerDecidedWrongMsgTypeDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingMultiSignerProposalMessageWithHeight(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						qbft.FirstHeight,
					),
				},
			},
		},
		nil,
		"could not process msg: invalid signed message: msg allows 1 signer",
		nil,
	)

	test.SetPrivateKeys(ks)

	return test
}
