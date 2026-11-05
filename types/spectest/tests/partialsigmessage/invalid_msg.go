package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidMsg tests a signed msg with 1 invalid message
func InvalidMsg() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, gloas.DataVersionGloas)
	msg.Messages = append(msg.Messages, &types.PartialSignatureMessage{})

	return NewMsgSpecTest(
		"invalid message",
		testdoc.MsgSpecTestInvalidMsgDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		types.InconsistentSignersErrorCode,
	)
}
