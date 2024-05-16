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

// SequencedHappyFlowDuties performs the happy flow of a sequence of duties
func SequencedHappyFlowDuties() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.CommitteeMsgID(ks)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "sequenced happy flow duties",
		Tests: []*committee.CommitteeSpecTest{},
	}

	// Returns a list of [Duty, qbft messages...] for each duty slot in sequence
	inputWithDecidingMessages := func(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []interface{} {

		attValidatorsIndexList := testingutils.ValidatorIndexList(numAttestingValidators)
		attKsMap := testingutils.KeySetMapForValidators(numAttestingValidators)

		scValidatorsIndexList := testingutils.ValidatorIndexList(numSyncCommitteeValidators)
		scKsMap := testingutils.KeySetMapForValidators(numSyncCommitteeValidators)

		jointMap := attKsMap
		for valIdx, valKS := range scKsMap {
			jointMap[valIdx] = valKS
		}

		ret := make([]interface{}, 0)
		for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

			currentSlot := testingutils.TestingDutySlot + slotIncrement

			// Duty
			duty := testingutils.TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList)
			ret = append(ret, duty)

			// QBFT
			for _, msg := range testingutils.SSVDecidingMsgsForHeightWithRoot(sha256.Sum256(testingutils.TestBeaconVoteByts), testingutils.TestBeaconVoteByts, msgID, qbft.Height(currentSlot), ks) {
				ret = append(ret, msg)
			}

			// Post-consensus
			for opID := uint64(1); opID <= ks.Threshold; opID++ {
				postConsensusMsg := testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(duty, jointMap, opID)))
				ret = append(ret, postConsensusMsg)
			}
		}
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

	beaconBroadcastedRoots := func(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []string {

		attValidatorsIndexList := testingutils.ValidatorIndexList(numAttestingValidators)
		attKsMap := testingutils.KeySetMapForValidators(numAttestingValidators)

		scValidatorsIndexList := testingutils.ValidatorIndexList(numSyncCommitteeValidators)
		scKsMap := testingutils.KeySetMapForValidators(numSyncCommitteeValidators)

		jointMap := attKsMap
		for valIdx, valKS := range scKsMap {
			jointMap[valIdx] = valKS
		}

		ret := make([]string, 0)
		for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

			currentSlot := testingutils.TestingDutySlot + slotIncrement

			duty := testingutils.TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList)

			ret = append(ret, testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(duty, jointMap)...)
		}
		return ret
	}

	for _, numSequencedDuties := range []int{1, 2, 4} {
		// TODO add 500
		for _, numValidators := range []int{1, 30} {

			ksMap := testingutils.KeySetMapForValidators(numValidators)
			shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

			multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
				{
					Name:                   fmt.Sprintf("%v duties %v attestation", numSequencedDuties, numValidators),
					Committee:              testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  inputWithDecidingMessages(numSequencedDuties, numValidators, 0),
					OutputMessages:         outputMessages(numSequencedDuties, numValidators, 0),
					BeaconBroadcastedRoots: beaconBroadcastedRoots(numSequencedDuties, numValidators, 0),
				},
				{
					Name:                   fmt.Sprintf("%v duties %v sync committee", numSequencedDuties, numValidators),
					Committee:              testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  inputWithDecidingMessages(numSequencedDuties, 0, numValidators),
					OutputMessages:         outputMessages(numSequencedDuties, 0, numValidators),
					BeaconBroadcastedRoots: beaconBroadcastedRoots(numSequencedDuties, 0, numValidators),
				},
				{
					Name:                   fmt.Sprintf("%v duties %v attestations %v sync committees", numSequencedDuties, numValidators, numValidators),
					Committee:              testingutils.BaseCommitteeWithRunnerSample(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
					Input:                  inputWithDecidingMessages(numSequencedDuties, numValidators, numValidators),
					OutputMessages:         outputMessages(numSequencedDuties, numValidators, numValidators),
					BeaconBroadcastedRoots: beaconBroadcastedRoots(numSequencedDuties, numValidators, numValidators),
				},
			}...)
		}
	}

	return multiSpecTest
}
