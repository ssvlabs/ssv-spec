package signedssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoSigner tests an invalid SignedSSVMessageTest with no signer
func NoSigner() *SignedSSVMessageTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SSVMsgAttester(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1))
	msg.OperatorID = []types.OperatorID{}

	return &SignedSSVMessageTest{
		Name:          "no signer",
		Messages:      []*types.SignedSSVMessage{msg},
		ExpectedError: "No OperatorID in SignedSSVMessage",
	}
}
