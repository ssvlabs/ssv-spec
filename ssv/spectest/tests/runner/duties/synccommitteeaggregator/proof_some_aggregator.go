package synccommitteeaggregator

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SomeAggregatorQuorum tests a quorum of selection proofs of which some are aggregator
func SomeAggregatorQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &SyncCommitteeAggregatorProofSpecTest{
		Name: "sync committee aggregator some are aggregators",
		Messages: []*types.SSVMessage{
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
		},
		ProofRootsMap: map[string]bool{
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[0][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[1][:]): false,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[2][:]): true,
		},
		PostDutyRunnerStateRoot: "609cb39abfee9a3308c679f5e88de8724215911d9e053c6972e00c32d94e2401",
	}
}
