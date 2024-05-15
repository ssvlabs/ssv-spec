package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapProposer tests mapping of BNRoleProposer.
func MapProposer() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map proposer role",
		BeaconRole: types.BNRoleProposer,
		RunnerRole: types.RoleProposer,
	}
}
