package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ZeroSigner tests an invalid SignedSSVMessageTest with zero signer
func ZeroSigner() *SignedSSVMessageTest {

	return &SignedSSVMessageTest{
		Name: "zero signer",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: 0,
				Signature:  testingutils.TestingSignedSSVMessageSignature,
				Data:       []byte{1, 2, 3, 4},
			},
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
