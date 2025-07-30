package duty

import (
	"math"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapUnknownRole tests mapping of an unknown role.
func MapUnknownRole() *DutySpecTest {
	return NewDutySpecTest(
		"map unknown role",
		testdoc.MapUnknownRoleTestDoc,
		math.MaxInt32,
		types.RoleUnknown,
	)
}
