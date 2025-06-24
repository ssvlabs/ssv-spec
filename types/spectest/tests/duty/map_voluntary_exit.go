package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapVoluntaryExit tests mapping of BNRoleVoluntaryExit.
func MapVoluntaryExit() *DutySpecTest {
	return NewDutySpecTest(
		"map voluntary exit role",
		types.BNRoleVoluntaryExit,
		types.RoleVoluntaryExit,
	)
}
