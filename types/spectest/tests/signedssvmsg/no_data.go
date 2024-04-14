package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NilSSVMessage tests an invalid SignedSSVMessageTest with nil SSVMessage
func NilSSVMessage() *SignedSSVMessageTest {

	return &SignedSSVMessageTest{
		Name: "nil ssvmessage",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: []types.OperatorID{1},
				Signature:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage: nil,
			},
		},
		ExpectedError: "nil SSVMessage",
	}
}
