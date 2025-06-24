package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapAggregator tests mapping of BNRoleAggregator.
func MapAggregator() *DutySpecTest {
	return NewDutySpecTest(
		"map aggregator role",
		types.BNRoleAggregator,
		types.RoleAggregator,
	)
}
