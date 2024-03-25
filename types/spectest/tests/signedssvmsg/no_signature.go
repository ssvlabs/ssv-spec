package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSignature tests an invalid SignedSSVMessageTest with no signature
func NoSignature() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	return &SignedSSVMessageTest{
		Name: "no signature",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: []types.OperatorID{1},
				Signature:  [][]byte{},
				SSVMessage: testingutils.SSVMsgAttester(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
			},
		},
		ExpectedError: "No signature in SignedSSVMessage",
	}
}
