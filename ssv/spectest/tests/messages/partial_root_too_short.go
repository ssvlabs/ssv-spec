package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// PartialRootTooShort tests PostConsensusMessage root < 32 bytes
func PartialRootTooShort() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages[0].SigningRoot = make([]byte, 31)

	return &MsgSpecTest{
		Name: "partial root too short",
		Messages: []*ssv.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "message invalid: SigningRoot invalid",
	}
}
