package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// validatorRegistrationBeaconRole is the reserved beacon role of the validator-registration duty.
const validatorRegistrationBeaconRole = types.BeaconRole(5)

// MapValidatorRegistration tests that the validator-registration duty maps to no runner: SIP #94 §5 deprecates
// it at Gloas in favour of proposer preferences, and its beacon role value stays reserved.
func MapValidatorRegistration() *DutySpecTest {
	return NewDutySpecTest(
		"map validator registration role",
		testdoc.MapValidatorRegistrationTestDoc,
		validatorRegistrationBeaconRole,
		types.RoleUnknown,
	)
}
