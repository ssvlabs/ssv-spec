package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSignature tests an invalid SignedSSVMessageTest with no signature
func NoSignature() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	msg.Signature = [][]byte{}

	return &SignedSSVMessageTest{
		Name:          "no signature",
		Messages:      []*types.SignedSSVMessage{msg},
		ExpectedError: "No signature in SignedSSVMessage",
	}
}
