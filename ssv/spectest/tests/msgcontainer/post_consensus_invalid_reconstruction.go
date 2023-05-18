package msgcontainer

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostConsensusInvalidReconstruction tests adding a single post consensus message to container and trying to reconstruct the full signature
func PostConsensusInvalidReconstruction() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "post consensus invalid reconstruction",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
		},
		PostMsgCount: 1,
		PostReconstructedSignature: []string{
			"",
			"",
			"",
		},
		ExpectedErr: "failed to verify reconstruct signature: could not reconstruct a valid signature",
	}
}
