package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ZeroSigner tests an invalid SignedSSVMessageTest with zero signer
func ZeroSigner() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	msg.OperatorID = []types.OperatorID{0}

	return &SignedSSVMessageTest{
		Name:          "zero signer",
		Messages:      []*types.SignedSSVMessage{msg},
		ExpectedError: "OperatorID in SignedSSVMessage is 0",
	}
}
