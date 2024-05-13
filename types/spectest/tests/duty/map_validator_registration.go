package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapValidatorRegistration tests mapping of BNRoleValidatorRegistration.
func MapValidatorRegistration() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map validator registration role",
		BeaconRole: types.BNRoleValidatorRegistration,
		RunnerRole: types.RoleValidatorRegistration,
	}
}
