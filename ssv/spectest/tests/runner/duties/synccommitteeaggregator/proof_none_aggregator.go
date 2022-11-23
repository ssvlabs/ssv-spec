package synccommitteeaggregator

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoneAggregatorQuorum tests a quorum of selection proofs but none of which are aggregator
func NoneAggregatorQuorum() *SyncCommitteeAggregatorProofSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &SyncCommitteeAggregatorProofSpecTest{
		Name: "sync committee aggregator none is aggregator",
		Messages: []*types.SSVMessage{
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
		},
		ProofRootsMap: map[string]bool{
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[0][:]): false,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[1][:]): false,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[2][:]): false,
		},
		PostDutyRunnerStateRoot: "bc86c14a5960cc85a51e50b5a8ad2fb72c85ffcae02a98892e33eb1d5877504b",
	}
}
