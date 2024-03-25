package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignerAndSignatureWithDifferentLength tests an invalid SignedSSVMessageTest with len(OperatorID) != len(Signature)
func SignerAndSignatureWithDifferentLength() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return &SignedSSVMessageTest{
		Name: "signer and signature with different length",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: []types.OperatorID{1, 2},
				Signature:  [][]byte{{1, 2, 3, 4}},
				SSVMessage: testingutils.SSVMsgAttester(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		ExpectedError: "SignedSSVMessage has a different number of operato IDs and signatures",
	}
}
