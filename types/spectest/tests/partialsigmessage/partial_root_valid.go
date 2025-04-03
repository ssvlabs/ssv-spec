package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialRootValid tests PostConsensusMessage root == 32 bytes
func PartialRootValid() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)

	return &MsgSpecTest{
		Name: "partial root valid",
		Messages: []*types.PartialSignatureMessages{
			msg,
		},
	}
}
