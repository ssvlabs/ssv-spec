package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
