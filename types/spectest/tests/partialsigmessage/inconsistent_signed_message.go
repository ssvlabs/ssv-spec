package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InconsistentSignedMessage tests SignedPartialSignatureMessage where the signer is not the same as the signer in messages
func InconsistentSignedMessage() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgWithDifferentSigner := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, spec.DataVersionPhase0)

	msg.Messages = append(msg.Messages, msgWithDifferentSigner.Messages...)

	return NewMsgSpecTest(
		"inconsistent signed message",
		testdoc.MsgSpecTestInconsistentSignedMessageDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		"inconsistent signers",
	)
}
