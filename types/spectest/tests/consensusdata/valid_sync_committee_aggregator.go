package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidSyncCommitteeAggregator tests a valid sync committee aggregator
func ValidSyncCommitteeAggregator() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "valid sync committee aggregator",
		Obj: &types.ConsensusData{
			Duty: testingutils.TestingSyncCommitteeContributionDuty,
			SyncCommitteeContribution: types.ContributionsMap{
				testingutils.TestingContributionProofsSigned[0]: testingutils.TestingSyncCommitteeContributions[0],
				testingutils.TestingContributionProofsSigned[1]: testingutils.TestingSyncCommitteeContributions[1],
				testingutils.TestingContributionProofsSigned[2]: testingutils.TestingSyncCommitteeContributions[2],
			},
		},
	}
}
