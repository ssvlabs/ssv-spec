package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapValidatorRegistration tests mapping of BNRoleValidatorRegistration.
func MapValidatorRegistration() *DutySpecTest {
	return NewDutySpecTest(
		"map validator registration role",
		"Test mapping of BNRoleValidatorRegistration",
		types.BNRoleValidatorRegistration,
		types.RoleValidatorRegistration,
	)
}
