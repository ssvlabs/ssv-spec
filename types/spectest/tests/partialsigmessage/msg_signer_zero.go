package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MessageSigner0 tests PartialSignatureMessage signer == 0
func MessageSigner0() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgPre := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgPre.Messages[0].Signer = 0
	msgPost := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msgPost.Messages[0].Signer = 0

	return NewMsgSpecTest(
		"message signer 0",
		testdoc.MsgSpecTestMessageSigner0Doc,
		[]*types.PartialSignatureMessages{msgPre, msgPost},
		nil,
		nil,
		"message invalid: signer ID 0 not allowed",
	)
}
