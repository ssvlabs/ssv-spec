package committeesingleduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// StartNoDuty starts a cluster runner with no duties
func StartNoDuty() tests.SpecTest {

	ksMapFor1Validator := testingutils.KeySetMapForValidators(1)

	return &committee.CommitteeSpecTest{
		Name:      "empty committee duty",
		Committee: testingutils.BaseCommittee(ksMapFor1Validator),
		Input: []interface{}{
			testingutils.TestingCommitteeDuty(nil, nil, spec.DataVersionElectra),
		},
		ExpectedError:  "no beacon duties",
		OutputMessages: []*types.PartialSignatureMessages{},
	}
}
