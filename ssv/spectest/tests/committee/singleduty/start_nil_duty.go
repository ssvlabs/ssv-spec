package committeesingleduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartNilDuty starts a nil committee duty
func StartNilDuty() tests.SpecTest {

	ksMapFor1Validator := testingutils.KeySetMapForValidators(1)

	return committee.NewCommitteeSpecTest(
		"nil committee duty",
		testdoc.CommitteeStartNilDutyDoc,
		testingutils.BaseCommittee(ksMapFor1Validator),
		[]interface{}{
			(*types.CommitteeDuty)(nil),
		},
		"",
		nil,
		nil,
		nil,
		types.UnknownDutyRoleDataErrorCode,
		nil,
	)
}
