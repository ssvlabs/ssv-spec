package committeemultipleduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ShuffledHappyFlowDutiesWithTheSameValidators performs the happy flow of duties with shuffled input messages (that preserves order between duty messages)
// The duties are assigned to the same validators. This causes the previous duties to be stopped and only the beacon root of the last duty is submitted
func ShuffledHappyFlowDutiesWithTheSameValidators() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "shuffled happy flow duties with same validators",
		Tests: []*committee.CommitteeSpecTest{},
	}

	expectedError := func(numOfDuties int) string {
		if numOfDuties == 1 {
			return ""
		}
		return "could not find validators for root"
	}

	for _, numSequencedDuties := range []int{1, 2, 4} {

		broadcastedBeaconRootSlot := phase0.Slot(testingutils.TestingDutySlot + numSequencedDuties - 1)

		for _, numValidators := range []int{1, 30, 500} {

			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
				{
					Name:                   fmt.Sprintf("%v duties %v attestation", numSequencedDuties, numValidators),
					Committee:              testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  testingutils.CommitteeInputForDutiesWithShuffle(numSequencedDuties, numValidators, 0, true),
					OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numSequencedDuties, numValidators, 0),
					BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDuty(broadcastedBeaconRootSlot, numValidators, 0),
					ExpectedError:          expectedError(numSequencedDuties),
				},
				{
					Name:                   fmt.Sprintf("%v duties %v sync committee", numSequencedDuties, numValidators),
					Committee:              testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  testingutils.CommitteeInputForDutiesWithShuffle(numSequencedDuties, 0, numValidators, true),
					OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numSequencedDuties, 0, numValidators),
					BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDuty(broadcastedBeaconRootSlot, 0, numValidators),
					ExpectedError:          expectedError(numSequencedDuties),
				},
				{
					Name:                   fmt.Sprintf("%v duties %v attestations %v sync committees", numSequencedDuties, numValidators, numValidators),
					Committee:              testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  testingutils.CommitteeInputForDutiesWithShuffle(numSequencedDuties, numValidators, numValidators, true),
					OutputMessages:         testingutils.CommitteeOutputMessagesForDuties(numSequencedDuties, numValidators, numValidators),
					BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDuty(broadcastedBeaconRootSlot, numValidators, numValidators),
					ExpectedError:          expectedError(numSequencedDuties),
				},
			}...)
		}
	}

	return multiSpecTest
}
