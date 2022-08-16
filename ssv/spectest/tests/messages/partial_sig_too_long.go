package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PartialSigTooLong tests PostConsensusMessage sig > 96 bytes
func PartialSigTooLong() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages[0].PartialSignature = make([]byte, 97)

	return &MsgSpecTest{
		Name: "partial sig too long",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "message invalid: PartialSignatureMessage sig invalid",
	}
}
