package committeesingleduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownValidator tests starting a duty for an unknown validator
func UnknownValidator() tests.SpecTest {

	committeeDutyExpectedErr := "no shares for duty's validators"
	expectedError := "unknown validator for duty"

	ksMap := testingutils.KeySetMapForValidators(10)
	unknownValidator := 11

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name: "unknown validator",
		Tests: []*committee.CommitteeSpecTest{
			{
				Name:      "sync committee",
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot, []int{unknownValidator}),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  committeeDutyExpectedErr,
			},
			{
				Name:      "attestation",
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, []int{unknownValidator}),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  committeeDutyExpectedErr,
			},
			{
				Name:      "proposer",
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingProposerDutyWithValidatorIndexV(phase0.ValidatorIndex(unknownValidator), spec.DataVersionDeneb),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      "aggregator",
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingAggregatorDutyWithValidatorIndex(phase0.ValidatorIndex(unknownValidator)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      "sync committee contribution",
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingSyncCommitteeContributionDutyWithValidatorIndex(phase0.ValidatorIndex(unknownValidator)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		},
	}
	return multiSpecTest
}
