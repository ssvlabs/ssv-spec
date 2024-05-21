package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"math"
)

// MapUnknownRole tests mapping of an unknown role.
func MapUnknownRole() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map unknown role",
		BeaconRole: math.MaxInt32,
		RunnerRole: types.RoleUnknown,
	}
}
