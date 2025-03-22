package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapProposer tests mapping of BNRoleProposer.
func MapCBSigning() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map commit boost signing role",
		BeaconRole: types.BNRoleCBSigning,
		RunnerRole: types.RoleCBSigning,
	}
}
