package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignersAndSignaturesWithDifferentLength tests an invalid SignedSSVMessageTest with len(signers) != len(signatures)
func SignersAndSignaturesWithDifferentLength() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()
	return NewSignedSSVMessageTest(
		"signers and signatures with different length",
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1, 2, 3, 4},
				Signatures:  [][]byte{{1, 2, 3, 4}, {2, 2, 3, 4}, {3, 2, 3, 4}},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		"number of signatures is different than number of signers",
		nil,
	)
}
