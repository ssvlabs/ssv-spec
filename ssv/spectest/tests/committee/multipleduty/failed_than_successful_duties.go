package committeemultipleduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FailedThanSuccessfulDuties decides a sequence of duties (not completing post-consensus) then performs the full happy flow of another sequence of duties
func FailedThanSuccessfulDuties() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "failed then successful duties",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		// TODO add 500
		for _, numValidators := range []int{1, 30} {
			for _, numFailedDuties := range []int{1, 2} {
				for _, numSuccessfulDuties := range []int{1, 2} {
					ksMap := testingutils.KeySetMapForValidators(numValidators)
					shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

					slot := testingutils.TestingDutySlotV(version) + phase0.Slot(numFailedDuties)

					multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
						{
							Name:                   fmt.Sprintf("%v fails %v successful for %v attestation (%s)", numFailedDuties, numSuccessfulDuties, numValidators, version.String()),
							Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
							Input:                  testingutils.CommitteeInputForDutiesWithFailuresThanSuccess(numValidators, 0, numFailedDuties, numSuccessfulDuties, version),
							OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numFailedDuties+numSuccessfulDuties, numValidators, 0, version),
							BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSuccessfulDuties, numValidators, 0, slot, version),
						},
						{
							Name:                   fmt.Sprintf("%v fails %v successful for %v sync committee (%s)", numFailedDuties, numSuccessfulDuties, numValidators, version.String()),
							Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
							Input:                  testingutils.CommitteeInputForDutiesWithFailuresThanSuccess(0, numValidators, numFailedDuties, numSuccessfulDuties, version),
							OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numFailedDuties+numSuccessfulDuties, 0, numValidators, version),
							BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSuccessfulDuties, 0, numValidators, slot, version),
						},
						{
							Name:                   fmt.Sprintf("%v fails %v successful for %v attestations %v sync committees (%s)", numFailedDuties, numSuccessfulDuties, numValidators, numValidators, version.String()),
							Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
							Input:                  testingutils.CommitteeInputForDutiesWithFailuresThanSuccess(numValidators, numValidators, numFailedDuties, numSuccessfulDuties, version),
							OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numFailedDuties+numSuccessfulDuties, numValidators, numValidators, version),
							BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSuccessfulDuties, numValidators, numValidators, slot, version),
						},
					}...)
				}
			}
		}
	}
	return multiSpecTest
}
