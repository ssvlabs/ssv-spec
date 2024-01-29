package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SigValid tests SignedPostConsensusMessage sig == 96 bytes
func SigValid() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)

	return &MsgSpecTest{
		Name: "sig valid",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
	}
}
