package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoMsgs tests a signed msg with no msgs
func NoMsgs() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, gloas.DataVersionGloas)
	msg.Messages = []*types.PartialSignatureMessage{}

	return NewMsgSpecTest(
		"no messages",
		testdoc.MsgSpecTestNoMsgsDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		types.NoPartialSigMessagesErrorCode,
	)
}
