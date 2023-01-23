package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// PartialSigTooShort tests PostConsensusMessage sig < 96 bytes
func PartialSigTooShort() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages[0].PartialSignature = make([]byte, 95)

	return &MsgSpecTest{
		Name: "partial sig too short",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "message invalid: PartialSignatureMessage sig invalid",
	}
}
