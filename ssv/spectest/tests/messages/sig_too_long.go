package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
