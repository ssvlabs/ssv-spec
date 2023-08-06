package msgcontainer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests adding a duplicate msg (same signer)
func DuplicateMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "duplicate msg",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionBellatrix),
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionBellatrix),
		},
		PostMsgCount:               1,
		PostReconstructedSignature: []string{},
	}
}
