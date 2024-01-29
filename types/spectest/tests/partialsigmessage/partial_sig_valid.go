package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PartialSigValid tests PostConsensusMessage sig == 96 bytes
func PartialSigValid() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)

	return &MsgSpecTest{
		Name: "partial sig valid",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
	}
}
