package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidMsg tests a signed msg with 1 invalid message
func InvalidMsg() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msg.Messages = append(msg.Messages, &types.PartialSignatureMessage{})

	return NewMsgSpecTest(
		"invalid message",
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		"inconsistent signers",
	)
}
