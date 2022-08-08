package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
