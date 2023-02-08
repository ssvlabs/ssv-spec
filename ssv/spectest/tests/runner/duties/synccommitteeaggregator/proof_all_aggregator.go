package synccommitteeaggregator

import (
	"encoding/hex"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
		PostDutyRunnerStateRoot: "d644d3960c5c038e1403483dcb9fabc72f84552f4f6996d83f71b7956af2ff13",
	}
}
