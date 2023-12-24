package msgcontainer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SinglePreConsensusMsg tests adding a single pre consensus message to container
func SinglePreConsensusMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "single pre consensus message",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionBellatrix),
		},
		PostMsgCount:               1,
		PostReconstructedSignature: []string{},
	}
}
