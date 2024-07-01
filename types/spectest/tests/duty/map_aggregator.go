package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapAggregator tests mapping of BNRoleAggregator.
func MapAggregator() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map aggregator role",
		BeaconRole: types.BNRoleAggregator,
		RunnerRole: types.RoleAggregator,
	}
}
