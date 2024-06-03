package committeemultipleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ShuffledHappyFlowDutiesWithDifferentValidators performs the happy flow of duties with shuffled input messages (that preserves order between duty messages)
// The duties are assigned to different validators
func ShuffledHappyFlowDutiesWithDifferentValidators() tests.SpecTest {

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "shuffled happy flow duties with different validators",
		Tests: []*committee.CommitteeSpecTest{},
	}

	// TODO add 500
	for _, numSequencedDuties := range []int{1, 2, 4} {
		for _, numValidators := range []int{8, 30} {

			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
				{
					Name:                   fmt.Sprintf("%v duties %v attestation", numSequencedDuties, numValidators),
					Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  testingutils.CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties, numValidators, 0, true),
					OutputMessages:         testingutils.CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties, numValidators, 0),
					BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties, numValidators, 0),
				},
				{
					Name:                   fmt.Sprintf("%v duties %v sync committee", numSequencedDuties, numValidators),
					Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  testingutils.CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties, 0, numValidators, true),
					OutputMessages:         testingutils.CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties, 0, numValidators),
					BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties, 0, numValidators),
				},
				{
					Name:                   fmt.Sprintf("%v duties %v attestations %v sync committees", numSequencedDuties, numValidators, numValidators),
					Committee:              testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  testingutils.CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties, numValidators, numValidators, true),
					OutputMessages:         testingutils.CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties, numValidators, numValidators),
					BeaconBroadcastedRoots: testingutils.CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties, numValidators, numValidators),
				},
			}...)
		}
	}

	return multiSpecTest
}
