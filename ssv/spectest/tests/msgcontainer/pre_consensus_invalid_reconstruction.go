package msgcontainer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreConsensusInvalidReconstruction tests adding a single pre consensus message to container and trying to reconstruct the full signature
func PreConsensusInvalidReconstruction() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "pre consensus invalid reconstruction",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionBellatrix),
		},
		PostMsgCount: 1,
		PostReconstructedSignature: []string{
			"",
		},
		ExpectedErr: "failed to verify reconstruct signature: could not reconstruct a valid signature",
	}
}
