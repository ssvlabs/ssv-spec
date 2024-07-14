package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidMsg tests a signed msg with 1 invalid message
func InvalidMsg() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Messages = append(msg.Messages, &types.PartialSignatureMessage{})

	return &MsgSpecTest{
		Name:          "invalid message",
		Messages:      []*types.PartialSignatureMessages{msg},
		ExpectedError: "inconsistent signers",
	}
}
