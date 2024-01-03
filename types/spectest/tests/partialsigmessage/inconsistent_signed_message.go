package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InconsistentSignedMessage tests SignedPartialSignatureMessage where the signer is not the same as the signer in messages
func InconsistentSignedMessage() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix)
	msgWithDifferentSigner := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)

	msg.Message.Messages = append(msg.Message.Messages, msgWithDifferentSigner.Message.Messages...)

	return &MsgSpecTest{
		Name: "message signer 0",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "inconsistent signers",
	}
}
