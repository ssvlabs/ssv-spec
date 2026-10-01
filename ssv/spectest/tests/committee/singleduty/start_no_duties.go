package committeesingleduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartNoDuty starts a cluster runner with no duties
func StartNoDuty() tests.SpecTest {

	ksMapFor1Validator := testingutils.KeySetMapForValidators(1)

	return committee.NewCommitteeSpecTest(
		"empty committee duty",
		testdoc.CommitteeStartNoDutyDoc,
		testingutils.BaseCommittee(ksMapFor1Validator),
		[]interface{}{
			testingutils.TestingCommitteeDuty(nil, nil, gloas.DataVersionGloas),
		},
		"",
		nil,
		nil,
		nil,
		types.NoBeaconDutiesErrorCode,
		nil,
	)
}
