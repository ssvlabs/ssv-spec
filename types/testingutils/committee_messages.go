package testingutils

import (
	"crypto/sha256"
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// This function return the validators that should be assigned to the duties
func selectValidatorsForDuties(numAttestingValidators int, numSyncCommitteeValidators int, numSequencedDuties int, currDutyIncrement int, diffValidatorsForDuties bool) ([]int, []int) {

	// Attestation validators
	attValidatorsIndexList := ValidatorIndexList(numAttestingValidators)

	// Sync committee validators
	scValidatorsIndexList := ValidatorIndexList(numSyncCommitteeValidators)

	// Return variables
	var attValidatorsForDuty []int
	var scValidatorsForDuty []int

	if diffValidatorsForDuties {
		// If not using the same validators, equally partition the validators list per duty
		attPartitionLength := len(attValidatorsIndexList) / numSequencedDuties
		attValidatorsForDuty = attValidatorsIndexList[currDutyIncrement*attPartitionLength : (currDutyIncrement+1)*attPartitionLength]
		scPartitionLength := len(scValidatorsIndexList) / numSequencedDuties
		scValidatorsForDuty = scValidatorsIndexList[currDutyIncrement*scPartitionLength : (currDutyIncrement+1)*scPartitionLength]
	} else {
		// If same validators, use all validators for the duties
		attValidatorsForDuty = attValidatorsIndexList
		scValidatorsForDuty = scValidatorsIndexList
	}
	return attValidatorsForDuty, scValidatorsForDuty
}

// Committee inputs

func CommitteeInputForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool) []interface{} {
	return CommitteeInputForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, addPostConsensus, false, false)
}

func CommitteeInputForDutiesWithShuffle(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool) []interface{} {
	return CommitteeInputForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, addPostConsensus, true, false)
}

func CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool) []interface{} {
	return CommitteeInputForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, addPostConsensus, true, true)
}

// Returns a list of [Duty, qbft messages...] for each duty
func CommitteeInputForDutiesWithParams(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool, shuffle bool, diffValidatorsForDuties bool) []interface{} {

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	// Sample KeySet
	var sampleKeySet *TestKeySet
	for _, ks := range ksMap {
		sampleKeySet = ks
		break
	}

	// Message ID
	msgID := CommitteeMsgID(sampleKeySet)

	ret := make([]interface{}, 0)
	dutiesCommands := make([]interface{}, 0)
	dutiesMsgs := make([][]interface{}, 0)

	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		dutyMsgs := make([]interface{}, 0)

		currentSlot := TestingDutySlot + slotIncrement

		// Validators assigned to the duties
		attValidatorsForDuty, scValidatorsForDuty := selectValidatorsForDuties(numAttestingValidators, numSyncCommitteeValidators, numSequencedDuties, slotIncrement, diffValidatorsForDuties)

		// Duty
		duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsForDuty, scValidatorsForDuty)
		dutiesCommands = append(dutiesCommands, duty)

		// QBFT
		for _, msg := range SSVDecidingMsgsForHeightWithRoot(sha256.Sum256(TestBeaconVoteByts), TestBeaconVoteByts, msgID, qbft.Height(currentSlot), sampleKeySet) {
			dutyMsgs = append(dutyMsgs, msg)
		}

		// Post-consensus
		if addPostConsensus {
			for opID := uint64(1); opID <= sampleKeySet.Threshold; opID++ {
				postConsensusMsg := SignPartialSigSSVMessage(sampleKeySet, SSVMsgCommittee(sampleKeySet, nil, PostConsensusCommitteeMsgForDuty(duty, ksMap, opID)))
				dutyMsgs = append(dutyMsgs, postConsensusMsg)
			}
		}

		dutiesMsgs = append(dutiesMsgs, dutyMsgs)
	}

	if shuffle {
		// If we should shuffle, insert duties command then shuffled duty msgs
		ret = append(ret, dutiesCommands...)
		ret = append(ret, MergeListsWithRandomPick(dutiesMsgs)...)
	} else {
		// If don't shuffle, insert duties command and msgs in order
		for i := 0; i < len(dutiesCommands); i++ {
			ret = append(ret, dutiesCommands[i])
			ret = append(ret, dutiesMsgs[i]...)
		}
	}

	return ret
}

// Committee output

func CommitteeOutputMessagesForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []*types.PartialSignatureMessages {
	return CommitteeOutputMessagesForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, false)
}

func CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []*types.PartialSignatureMessages {
	return CommitteeOutputMessagesForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, true)
}

// Returns a list of partial sig output messages for each duty slot
func CommitteeOutputMessagesForDutiesWithParams(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, diffValidatorsForDuties bool) []*types.PartialSignatureMessages {

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	ret := make([]*types.PartialSignatureMessages, 0)

	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		currentSlot := TestingDutySlot + slotIncrement

		// Validators assigned to the duties
		attValidatorsForDuty, scValidatorsForDuty := selectValidatorsForDuties(numAttestingValidators, numSyncCommitteeValidators, numSequencedDuties, slotIncrement, diffValidatorsForDuties)

		duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsForDuty, scValidatorsForDuty)
		postConsensusMsg := PostConsensusCommitteeMsgForDuty(duty, ksMap, 1)

		// Add post consensus
		ret = append(ret, postConsensusMsg)
	}

	return ret
}

// Committee beacon roots

func CommitteeBeaconBroadcastedRootsForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []string {
	return CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, false)
}
func CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []string {
	return CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, true)
}

// Returns a list of beacon roots for committees duties
func CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, diffValidatorsForDuties bool) []string {

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	ret := make([]string, 0)
	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		currentSlot := TestingDutySlot + slotIncrement

		// Validators assigned to the duties
		attValidatorsForDuty, scValidatorsForDuty := selectValidatorsForDuties(numAttestingValidators, numSyncCommitteeValidators, numSequencedDuties, slotIncrement, diffValidatorsForDuties)

		duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsForDuty, scValidatorsForDuty)

		beaconRoots := TestingSignedCommitteeBeaconObjectSSZRoot(duty, ksMap)

		// Add beacon roots
		ret = append(ret, beaconRoots...)
	}
	return ret
}

// Returns a list of beacon roots for committees duties
func CommitteeBeaconBroadcastedRootsForDuty(slot phase0.Slot, numAttestingValidators int, numSyncCommitteeValidators int) []string {

	// Attestation validators and KeySet map
	attValidatorsIndexList := ValidatorIndexList(numAttestingValidators)

	// Sync committee validators and KeySet map
	scValidatorsIndexList := ValidatorIndexList(numSyncCommitteeValidators)

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	ret := make([]string, 0)

	duty := TestingCommitteeDuty(slot, attValidatorsIndexList, scValidatorsIndexList)

	beaconRoots := TestingSignedCommitteeBeaconObjectSSZRoot(duty, ksMap)

	// Add beacon roots
	ret = append(ret, beaconRoots...)

	return ret
}
