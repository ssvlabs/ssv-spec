package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapAttester tests mapping of BNRoleAttester.
func MapAttester() *DutySpecTest {
	return NewDutySpecTest(
		"map attester role",
		testdoc.MapAttesterTestDoc,
		types.BNRoleAttester,
		types.RoleCommittee,
	)
}
