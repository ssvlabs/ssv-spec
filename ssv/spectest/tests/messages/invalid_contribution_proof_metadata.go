package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidContributionProofMetaData tests a PartialSignatureMessage for contribution proof metadata nil
func InvalidContributionProofMetaData() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Type = ssv.ContributionProofs
	msg.Message.Messages[0].MetaData = nil

	return &MsgSpecTest{
		Name:          "invalid meta data when type ContributionProofs",
		Messages:      []*ssv.SignedPartialSignatureMessage{msg},
		ExpectedError: "metadata nil for contribution proofs",
	}
}
