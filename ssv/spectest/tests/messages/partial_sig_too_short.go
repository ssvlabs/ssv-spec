package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PartialSigTooShort tests PostConsensusMessage sig < 96 bytes
func PartialSigTooShort() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.PartialSignatures[0].Signature = make([]byte, 95)

	return &MsgSpecTest{
		Name: "partial sig too short",
		Messages: []*ssv.SignedPartialSignatures{
			msg,
		},
		ExpectedError: "message invalid: PartialSignature sig invalid",
	}
}
