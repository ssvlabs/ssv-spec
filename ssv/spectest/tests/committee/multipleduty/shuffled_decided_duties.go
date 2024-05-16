package committeemultipleduty

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ShuffledDecidedDuties decides duties with shuffled input messages (that preserves order between duty messages)
func ShuffledDecidedDuties() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.CommitteeMsgID(ks)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "shuffled decided duties",
		Tests: []*committee.CommitteeSpecTest{},
	}

	// Returns a shuffled list (that preserves order between duty messages) of [Duty, qbft messages...] for each duty slot in sequence
	inputWithDecidingMessages := func(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []interface{} {

		attValidatorsIndexList := testingutils.ValidatorIndexList(numAttestingValidators)
		scValidatorsIndexList := testingutils.ValidatorIndexList(numSyncCommitteeValidators)

		ret := make([]interface{}, 0)
		dutiesMsgs := make([][]interface{}, 0)
		for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

			dutyMsgs := make([]interface{}, 0)

			currentSlot := testingutils.TestingDutySlot + slotIncrement

			// Set duties to ret in fixed slot order to guarantee that each is executed
			ret = append(ret, testingutils.TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList))

			for _, msg := range testingutils.SSVDecidingMsgsForHeightWithRoot(sha256.Sum256(testingutils.TestBeaconVoteByts), testingutils.TestBeaconVoteByts, msgID, qbft.Height(currentSlot), ks) {
				dutyMsgs = append(dutyMsgs, msg)
			}

			dutiesMsgs = append(dutiesMsgs, dutyMsgs)
		}

		ret = append(ret, testingutils.MergeListsWithRandomPick(dutiesMsgs)...)

		return ret
	}

	// Returns a list of output messages for each duty slot in sequence
	outputMessages := func(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []*types.PartialSignatureMessages {

		attValidatorsIndexList := testingutils.ValidatorIndexList(numAttestingValidators)
		attKsMap := testingutils.KeySetMapForValidators(numAttestingValidators)

		scValidatorsIndexList := testingutils.ValidatorIndexList(numSyncCommitteeValidators)
		scKsMap := testingutils.KeySetMapForValidators(numSyncCommitteeValidators)

		jointMap := attKsMap
		for valIdx, valKS := range scKsMap {
			jointMap[valIdx] = valKS
		}

		ret := make([]*types.PartialSignatureMessages, 0)
		for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

			currentSlot := testingutils.TestingDutySlot + slotIncrement

			duty := testingutils.TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList)

			postConsensusMsg := testingutils.PostConsensusCommitteeMsgForDuty(duty, jointMap, 1)
			if postConsensusMsg == nil {
				panic("post consensus message is nil")
			}
			ret = append(ret, postConsensusMsg)
		}
		return ret
	}

	for _, numSequencedDuties := range []int{2, 4} {
		// TODO add 500
		for _, numValidators := range []int{1, 30} {

			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
				{
					Name:           fmt.Sprintf("%v duties %v attestation", numSequencedDuties, numValidators),
					Committee:      testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:          inputWithDecidingMessages(numSequencedDuties, numValidators, 0),
					OutputMessages: outputMessages(numSequencedDuties, numValidators, 0),
				},
				{
					Name:           fmt.Sprintf("%v duties %v sync committee", numSequencedDuties, numValidators),
					Committee:      testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:          inputWithDecidingMessages(numSequencedDuties, 0, numValidators),
					OutputMessages: outputMessages(numSequencedDuties, 0, numValidators),
				},
				{
					Name:           fmt.Sprintf("%v duties %v attestations %v sync committees", numSequencedDuties, numValidators, numValidators),
					Committee:      testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:          inputWithDecidingMessages(numSequencedDuties, numValidators, numValidators),
					OutputMessages: outputMessages(numSequencedDuties, numValidators, numValidators),
				},
			}...)
		}
	}

	return multiSpecTest
}
