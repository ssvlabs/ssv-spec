package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NonUniqueSigner tests an invalid SignedSSVMessageTest with non unique signers
func NonUniqueSigner() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return NewSignedSSVMessageTest(
		"non unique signers",
		[]*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{1, 2, 2},
				Signatures:  [][]byte{{1, 2, 3, 4}, {2, 3, 4, 5}, {2, 3, 4, 5}},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		"non unique signer",
		nil,
	)
}
