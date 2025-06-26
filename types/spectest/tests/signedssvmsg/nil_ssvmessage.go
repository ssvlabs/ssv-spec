package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NilSSVMessage tests an invalid SignedSSVMessageTest with nil SSVMessage
func NilSSVMessage() *SignedSSVMessageTest {

	return NewSignedSSVMessageTest(
		"nil ssvmessage",
		"Test validation error for signed SSV message with nil SSVMessage",
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1},
				Signatures:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage:  nil,
			},
		},
		"nil SSVMessage",
		nil,
	)
}
