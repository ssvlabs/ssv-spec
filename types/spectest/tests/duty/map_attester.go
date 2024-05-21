package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapAttester tests mapping of BNRoleAttester.
func MapAttester() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map attester role",
		BeaconRole: types.BNRoleAttester,
		RunnerRole: types.RoleCommittee,
	}
}
