package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// UnknownType tests a validator duty whose beacon role doesn't map to a runner role
func UnknownType() *ValidationTest {
	duty := testingAttesterDuty()
	duty.Type = types.BNRoleUnknown

	return NewValidationTest(
		"unknown type",
		testdoc.ValidatorDutyValidationUnknownTypeDoc,
		[]*types.ValidatorDuty{duty},
		types.InvalidValidatorDutyErrorCode,
	)
}
