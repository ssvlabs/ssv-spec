package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Aggregator Committee Duty - Legacy functions for backward compatibility
// ==================================================

// TestingAggregatorCommitteeDutyOnlyAggregator creates a duty with only aggregator validators
func TestingAggregatorCommitteeDutyOnlyAggregator(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDutySingle(version)
}

// TestingAggregatorCommitteeDutyOnlySyncCommittee creates a duty with only sync committee validators
func TestingAggregatorCommitteeDutyOnlySyncCommittee() *types.AggregatorCommitteeDuty {
	return TestingSyncCommitteeContributionDuty
}

// TestingAggregatorCommitteeDutyMixed creates a duty with both aggregator and sync committee validators
func TestingAggregatorCommitteeDutyMixed(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorAndSyncCommitteeContributorDuties(version)
}

// TestingAggregatorCommitteeDutyMultipleAggregators creates a duty with multiple aggregator validators
func TestingAggregatorCommitteeDutyMultipleAggregators(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorDutyForValidators(version, []int{1, 2})
}

// TestingAggregatorCommitteeDutyMultipleSyncCommittee creates a duty with multiple sync committee validators
func TestingAggregatorCommitteeDutyMultipleSyncCommittee() *types.AggregatorCommitteeDuty {
	return TestingSyncCommitteeContributionDutyForValidators([]int{1, 2})
}
