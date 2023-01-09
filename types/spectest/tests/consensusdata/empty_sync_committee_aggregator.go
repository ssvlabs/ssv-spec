package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// EmptySyncCommitteeAggregator tests a invalid sync committee aggregator
func EmptySyncCommitteeAggregator() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "empty sync committee aggregator",
		Obj: &types.ConsensusData{
			Duty:                      testingutils.TestingSyncCommitteeContributionDuty,
			SyncCommitteeContribution: types.ContributionsMap{},
		},
		ExpectedErr: "sync committee contribution data is nil",
	}
}
