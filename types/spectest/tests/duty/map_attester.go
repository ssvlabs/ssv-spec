package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapAttester tests mapping of BNRoleAttester.
func MapAttester() *DutySpecTest {
	return NewDutySpecTest(
		"map attester role",
		types.BNRoleAttester,
		types.RoleCommittee,
	)
}
