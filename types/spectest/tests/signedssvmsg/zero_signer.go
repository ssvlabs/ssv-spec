package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ZeroSigner tests an invalid SignedSSVMessageTest with zero signer
func ZeroSigner() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return &SignedSSVMessageTest{
		Name: "zero signer",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: []types.OperatorID{0},
				Signature:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage: testingutils.SSVMsgAttester(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		ExpectedError: "OperatorID in SignedSSVMessage is 0",
	}
}
