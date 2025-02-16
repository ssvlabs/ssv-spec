package committeesingleduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartWithNoSharesForDuty starts a committee runner for a duty for which it doesn't have shares
func StartWithNoSharesForDuty() tests.SpecTest {

	// KeyShare map with entry only for validator 1
	ksMapFor1Validator := testingutils.KeySetMapForValidators(1)

	return &committee.CommitteeSpecTest{
		Name:      "start with no shares for duty",
		Committee: testingutils.BaseCommittee(ksMapFor1Validator),
		Input: []interface{}{
			// Duty for validator of index 2
			testingutils.TestingAttesterDutyForValidators(spec.DataVersionElectra, []int{2}),
		},
		ExpectedError:  "no shares for duty's validators",
		OutputMessages: []*types.PartialSignatureMessages{},
	}
}
