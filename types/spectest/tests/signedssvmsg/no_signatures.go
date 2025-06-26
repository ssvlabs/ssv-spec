package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoSignatures tests an invalid SignedSSVMessageTest with no signatures
func NoSignatures() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return NewSignedSSVMessageTest(
		"no signatures",
		"Test validation error for signed SSV message with no signatures",
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1},
				Signatures:  [][]byte{},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		"no signatures",
		nil,
	)
}
