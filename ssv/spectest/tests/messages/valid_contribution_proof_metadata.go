package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidContributionProofMetaData tests a PartialSignatureMessage for contribution proof metadata valid
func ValidContributionProofMetaData() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Type = ssv.ContributionProofs
	msg.Message.Messages[0].MetaData = &ssv.PartialSignatureMetaData{
		ContributionSubCommitteeIndex: 1,
	}

	return &MsgSpecTest{
		Name:     "valid meta data when type ContributionProofs",
		Messages: []*ssv.SignedPartialSignatureMessage{msg},
	}
}
