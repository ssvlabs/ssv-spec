package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidMsg tests a signed msg with 1 invalid message
func InvalidMsg() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages = append(msg.Message.Messages, &ssv.PartialSignatureMessage{})

	return &MsgSpecTest{
		Name:          "invalid message",
		Messages:      []*ssv.SignedPartialSignatureMessage{msg},
		ExpectedError: "inconsistent signers",
	}
}
