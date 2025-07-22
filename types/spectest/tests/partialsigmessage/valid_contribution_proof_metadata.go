package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidContributionProofMetaData tests a PartialSignatureMessage for contribution proof metadata valid
func ValidContributionProofMetaData() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, spec.DataVersionPhase0)
	msg.Type = types.ContributionProofs

	return NewMsgSpecTest(
		"valid meta data when type ContributionProofs",
		testdoc.MsgSpecTestValidContributionProofMetaDataDoc,
		[]*types.PartialSignatureMessages{msg},
		nil,
		nil,
		"",
	)
}
