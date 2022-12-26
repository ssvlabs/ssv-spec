package synccommitteeaggregator

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AllAggregatorQuorum tests a quorum of selection proofs of which all are aggregator
func AllAggregatorQuorum() *SyncCommitteeAggregatorProofSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &SyncCommitteeAggregatorProofSpecTest{
		Name: "sync committee aggregator all are aggregators",
		Messages: []*types.SSVMessage{
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
		},
		ProofRootsMap: map[string]bool{
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[0][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[1][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[2][:]): true,
		},
		PostDutyRunnerStateRoot: "03a9683e5a3c172b7b6e5fb2d3d63c2f7a7d4bbb4b5301b1f9d28f1e71b8fe39",
	}
}
