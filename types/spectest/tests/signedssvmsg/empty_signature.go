package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EmptySignature tests an invalid SignedSSVMessageTest with empty signature
func EmptySignature() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return NewSignedSSVMessageTest(
		"empty signature",
		"Test validation error for signed SSV message with empty signature",
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1},
				Signatures:  [][]byte{{}},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		"empty signature",
		nil,
	)
}
