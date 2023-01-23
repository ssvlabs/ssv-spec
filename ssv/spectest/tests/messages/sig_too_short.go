package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SigTooShort tests SignedPostConsensusMessage sig < 96 bytes
func SigTooShort() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Signature = make([]byte, 95)

	return &MsgSpecTest{
		Name: "sig too short",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "SignedPartialSignatureMessage sig invalid",
	}
}
