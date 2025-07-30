package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapVoluntaryExit tests mapping of BNRoleVoluntaryExit.
func MapVoluntaryExit() *DutySpecTest {
	return NewDutySpecTest(
		"map voluntary exit role",
		testdoc.MapVoluntaryExitTestDoc,
		types.BNRoleVoluntaryExit,
		types.RoleVoluntaryExit,
	)
}
