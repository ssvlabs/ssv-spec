package committeemultipleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FailedThanSuccessfulDuties decides a sequence of duties (not completing post-consensus) then performs the full happy flow of another sequence of duties
func FailedThanSuccessfulDuties() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "failed than successful duties",
		Tests: []*committee.CommitteeSpecTest{},
	}

	// TODO add 500
	for _, numValidators := range []int{1, 30} {
		for _, numFailedDuties := range []int{1, 2} {
			for _, numSuccessfulDuties := range []int{1, 2} {
				ksMap := testingutils.KeySetMapForValidators(numValidators)
				shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

				multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
					{
						Name:                   fmt.Sprintf("%v fails %v sucessful for%v attestation", numFailedDuties, numSuccessfulDuties, numValidators),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDutiesWithFailuresThanSuccess(numValidators, 0, numFailedDuties, numSuccessfulDuties),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numFailedDuties+numSuccessfulDuties, numValidators, 0),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSuccessfulDuties, numValidators, 0, testingutils.TestingDutySlot+numFailedDuties),
					},
					{
						Name:                   fmt.Sprintf("%v fails %v sucessful for %v sync committee", numFailedDuties, numSuccessfulDuties, numValidators),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDutiesWithFailuresThanSuccess(0, numValidators, numFailedDuties, numSuccessfulDuties),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numFailedDuties+numSuccessfulDuties, 0, numValidators),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSuccessfulDuties, 0, numValidators, testingutils.TestingDutySlot+numFailedDuties),
					},
					{
						Name:                   fmt.Sprintf("%v fails %v sucessful for %v attestations %v sync committees", numFailedDuties, numSuccessfulDuties, numValidators, numValidators),
						Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
						Input:                  testingutils.CommitteeInputForDutiesWithFailuresThanSuccess(numValidators, numValidators, numFailedDuties, numSuccessfulDuties),
						OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numFailedDuties+numSuccessfulDuties, numValidators, numValidators),
						BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSuccessfulDuties, numValidators, numValidators, testingutils.TestingDutySlot+numFailedDuties),
					},
				}...)
			}
		}
	}
	return multiSpecTest
}
