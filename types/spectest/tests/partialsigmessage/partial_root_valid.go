package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PartialRootValid tests PreConsensusMessage root == 32 bytes
func PartialRootValid() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)

	return &MsgSpecTest{
		Name: "partial root valid",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
	}
}
