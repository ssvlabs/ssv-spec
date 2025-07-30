package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapValidatorRegistration tests mapping of BNRoleValidatorRegistration.
func MapValidatorRegistration() *DutySpecTest {
	return NewDutySpecTest(
		"map validator registration role",
		testdoc.MapValidatorRegistrationTestDoc,
		types.BNRoleValidatorRegistration,
		types.RoleValidatorRegistration,
	)
}
