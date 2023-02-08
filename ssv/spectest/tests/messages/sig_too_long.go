package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SigTooLong tests SignedPostConsensusMessage sig > 96 bytes
func SigTooLong() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Signature = make([]byte, 97)

	return &MsgSpecTest{
		Name: "sig too long",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "SignedPartialSignatureMessage sig invalid",
	}
}
