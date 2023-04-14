package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PartialRootValid tests PostConsensusMessage root == 32 bytes
func PartialRootValid() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)

	return &MsgSpecTest{
		Name: "partial root valid",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
	}
}
