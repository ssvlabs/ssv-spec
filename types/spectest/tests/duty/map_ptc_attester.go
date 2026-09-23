package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapPTCAttester tests mapping of BNRolePTCAttester.
func MapPTCAttester() *DutySpecTest {
	return NewDutySpecTest(
		"map ptc attester role",
		testdoc.MapPTCAttesterTestDoc,
		types.BNRolePTCAttester,
		types.RolePTCAttester,
	)
}
