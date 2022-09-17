package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMsg tests a signed msg with 1 invalid message
func InvalidMsg() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages = append(msg.Message.Messages, &ssv.PartialSignature{})

	return &MsgSpecTest{
		Name:          "invalid message",
		Messages:      []*ssv.SignedPartialSignature{msg},
		ExpectedError: "message invalid: PartialSignature sig invalid",
	}
}
