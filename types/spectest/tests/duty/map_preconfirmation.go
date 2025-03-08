package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapProposer tests mapping of BNRoleProposer.
func MapPreconf() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map preconf role",
		BeaconRole: types.BNRolePreconfirmation,
		RunnerRole: types.RolePreconfirmation,
	}
}
