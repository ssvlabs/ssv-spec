package msgcontainer

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SinglePostConsensusMsg tests adding a single post consensus message to container
func SinglePostConsensusMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "single post consensus message",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
		},
		PostMsgCount:               1,
		PostReconstructedSignature: []string{},
	}
}
