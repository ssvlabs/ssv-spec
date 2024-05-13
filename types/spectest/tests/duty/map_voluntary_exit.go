package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapVoluntaryExit tests mapping of BNRoleVoluntaryExit.
func MapVoluntaryExit() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map voluntary exit role",
		BeaconRole: types.BNRoleVoluntaryExit,
		RunnerRole: types.RoleVoluntaryExit,
	}
}
