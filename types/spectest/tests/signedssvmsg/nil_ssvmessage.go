package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NilSSVMessage tests an invalid SignedSSVMessageTest with nil SSVMessage
func NilSSVMessage() *SignedSSVMessageTest {

	return &SignedSSVMessageTest{
		Name: "nil ssvmessage",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1},
				Signatures:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage:  nil,
			},
		},
		ExpectedError: "nil SSVMessage",
	}
}
