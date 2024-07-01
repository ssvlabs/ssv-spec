package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ZeroSigner tests an invalid SignedSSVMessageTest with zero signer
func ZeroSigner() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return &SignedSSVMessageTest{
		Name: "zero signer",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorIDs: []types.OperatorID{0},
				Signatures:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage:  testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
