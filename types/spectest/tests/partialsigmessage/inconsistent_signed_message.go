package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InconsistentSignedMessage tests SignedPartialSignatureMessage where the signer is not the same as the signer in messages
func InconsistentSignedMessage() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgWithDifferentSigner := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)

	msg.Messages = append(msg.Messages, msgWithDifferentSigner.Messages...)

	return &MsgSpecTest{
		Name: "inconsistent signed message",
		Messages: []*types.PartialSignatureMessages{
			msg,
		},
		ExpectedError: "inconsistent signers",
	}
}
