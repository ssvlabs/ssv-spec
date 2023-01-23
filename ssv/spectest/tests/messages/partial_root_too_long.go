package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// PartialRootTooLong tests PostConsensusMessage root > 32 bytes
func PartialRootTooLong() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages[0].SigningRoot = make([]byte, 33)

	return &MsgSpecTest{
		Name: "partial root too long",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "message invalid: SigningRoot invalid",
	}
}
