package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InconsistentSignedMessage tests SignedPartialSignatureMessage where the signer is not the same as the signer in messages
func InconsistentSignedMessage() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, gloas.DataVersionGloas)
	msgWithDifferentSigner := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, gloas.DataVersionGloas)

	msg.Messages = append(msg.Messages, msgWithDifferentSigner.Messages...)

	return NewMsgSpecTest(
		"inconsistent signed message",
		testdoc.MsgSpecTestInconsistentSignedMessageDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		types.InconsistentSignersErrorCode,
	)
}
