package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// Valid tests a valid attester validator duty
func Valid() *ValidationTest {
	return NewValidationTest(
		"valid",
		testdoc.ValidatorDutyValidationValidDoc,
		[]*types.ValidatorDuty{testingAttesterDuty()},
		0,
	)
}
