package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSigner tests an invalid SignedSSVMessageTest with no signer
func NoSigner() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return &SignedSSVMessageTest{
		Name: "no signer",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: []types.OperatorID{},
				Signature:  [][]byte{testingutils.TestingSignedSSVMessageSignature},
				SSVMessage: testingutils.SSVMsgAttester(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		ExpectedError: "No OperatorID in SignedSSVMessage",
	}
}
