package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapAggregator tests mapping of BNRoleAggregator.
func MapAggregator() *DutySpecTest {
	return NewDutySpecTest(
		"map aggregator role",
		testdoc.MapAggregatorTestDoc,
		types.BNRoleAggregator,
		types.RoleAggregator,
	)
}
