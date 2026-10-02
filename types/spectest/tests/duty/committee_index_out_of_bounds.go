package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// CommitteeIndexOutOfBounds tests an attester duty whose ValidatorCommitteeIndex is not below CommitteeLength
func CommitteeIndexOutOfBounds() *ValidationTest {
	duty := testingAttesterDuty()
	duty.ValidatorCommitteeIndex = duty.CommitteeLength

	return NewValidationTest(
		"committee index out of bounds",
		testdoc.ValidatorDutyValidationCommitteeIndexOutOfBoundsDoc,
		[]*types.ValidatorDuty{duty},
		types.InvalidValidatorDutyErrorCode,
	)
}
