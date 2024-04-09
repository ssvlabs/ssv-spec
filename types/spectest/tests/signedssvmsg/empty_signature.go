package signedssvmsg

import "github.com/bloxapp/ssv-spec/types"

// EmptySignature tests an invalid SignedSSVMessageTest with empty signature
func EmptySignature() *SignedSSVMessageTest {

	return &SignedSSVMessageTest{
		Name: "empty signature",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: 1,
				Signature:  [256]byte{},
				Data:       []byte{1, 2, 3, 4},
			},
		},
		ExpectedError: "could not decode SSVMessage from data in SignedSSVMessage: incorrect size",
	}
}
