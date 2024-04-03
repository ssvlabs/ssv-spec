package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// EmptySignature tests an invalid SignedSSVMessageTest with empty signature
func EmptySignature() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	msg.Signature = [][]byte{{}}

	return &SignedSSVMessageTest{
		Name:          "empty signature",
		Messages:      []*types.SignedSSVMessage{msg},
		ExpectedError: "Signature has length 0 in SignedSSVMessage",
	}
}
