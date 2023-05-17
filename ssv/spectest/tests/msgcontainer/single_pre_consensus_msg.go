package msgcontainer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// QuorumPreConsensusMsg tests adding a quorum of pre consensus message to container
func QuorumPreConsensusMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &SpecTest{
		Name: "single pre consensus message",
		MsgsToAdd: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionBellatrix),
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionBellatrix),
			testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionBellatrix),
		},
		PostMsgCount: 3,
		PostReconstructedSignature: []string{
			"b930ac1fafd0ab1c5dab9606ace35788afc4008acf867a7346d192af233d3f54e95ddb9b248d38efcb3edc1ae5fe571008846769e0aa4cdb5d8ce480efa55e5e74e40812e4f8c3f62f89cae7936b42c5e270e3dec5b2c559463e7319a1c5d16f",
		},
	}
}
