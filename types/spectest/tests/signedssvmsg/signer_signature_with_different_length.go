package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignerAndSignatureWithDifferentLength tests an invalid SignedSSVMessageTest with len(OperatorID) != len(Signature)
func SignerAndSignatureWithDifferentLength() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	msg.OperatorID = []types.OperatorID{1, 2}

	return &SignedSSVMessageTest{
		Name:          "signer and signature with different length",
		Messages:      []*types.SignedSSVMessage{msg},
		ExpectedError: "SignedSSVMessage has a different number of operato IDs and signatures",
	}
}
