package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InconsistentSignedMessage tests SignedPartialSignatureMessage where the signer is not the same as the signer in messages
func InconsistentSignedMessage() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgWithDifferentSigner := testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)

	msg.Message.Messages = append(msg.Message.Messages, msgWithDifferentSigner.Message.Messages...)

	return &MsgSpecTest{
		Name: "inconsistent signed message",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "inconsistent signers",
	}
}
