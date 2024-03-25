package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SSVMessageNil tests a SignedSSVMessageTest with a nil SSVMessage
func SSVMessageNil() *SignedSSVMessageTest {
	return &SignedSSVMessageTest{
		Name: "ssvmessage nil",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: []types.OperatorID{1},
				Signature:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage: nil,
			},
		},
		ExpectedError: "SSVMessage is nil",
	}
}
