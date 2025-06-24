package duty

import (
	"math"

	"github.com/ssvlabs/ssv-spec/types"
)

// MapUnknownRole tests mapping of an unknown role.
func MapUnknownRole() *DutySpecTest {
	return NewDutySpecTest(
		"map unknown role",
		math.MaxInt32,
		types.RoleUnknown,
	)
}
