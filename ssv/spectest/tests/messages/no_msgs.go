package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoMsgs tests a signed msg with no msgs
func NoMsgs() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages = []*ssv.PartialSignatureMessage{}

	return &MsgSpecTest{
		Name:          "no messages",
		Messages:      []*ssv.SignedPartialSignatureMessage{msg},
		ExpectedError: "no PartialSignatureMessages messages",
	}
}
