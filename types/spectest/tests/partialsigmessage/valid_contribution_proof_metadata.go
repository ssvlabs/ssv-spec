package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidContributionProofMetaData tests a PartialSignatureMessage for contribution proof metadata valid
func ValidContributionProofMetaData() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)
	msg.Message.Type = types.ContributionProofs

	return &MsgSpecTest{
		Name:     "valid meta data when type ContributionProofs",
		Messages: []*types.SignedPartialSignatureMessage{msg},
	}
}
