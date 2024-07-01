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
	return CommitteeInputForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, addPostConsensus, false, false)
}

func CommitteeInputForDutiesWithShuffle(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool) []interface{} {
	return CommitteeInputForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, addPostConsensus, true, false)
}

func CommitteeInputForDutiesWithShuffleAndDifferentValidators(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool) []interface{} {
	return CommitteeInputForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, addPostConsensus, true, true)
}

func CommitteeInputForDutiesWithFailuresThanSuccess(numAttestingValidators int, numSyncCommitteeValidators int, numFails int, numSuccess int) []interface{} {
	ret := CommitteeInputForDutiesWithParams(numFails, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, false, false, false)
	ret = append(ret, CommitteeInputForDutiesWithParams(numSuccess, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot+numFails, true, false, false)...)
	return ret
}

// Returns a list of [Duty, qbft messages...] for each duty
func CommitteeInputForDutiesWithParams(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, startingSlot int, addPostConsensus bool, shuffle bool, diffValidatorsForDuties bool) []interface{} {

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	ret := make([]interface{}, 0)
	dutiesCommands := make([]interface{}, 0)
	dutiesMsgs := make([][]interface{}, 0)

	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		// Current slot
		currentSlot := startingSlot + slotIncrement

		// Duty and messages for the given slot
		duty, msgs := CommitteeInputForSlotInSequencedDuties(numAttestingValidators, numSyncCommitteeValidators, numSequencedDuties, slotIncrement, currentSlot, ksMap, addPostConsensus, diffValidatorsForDuties)

		// Add duty to list
		dutiesCommands = append(dutiesCommands, duty)

		// Add duty's messages list to dutiesMsgs
		dutyMsgs := make([]interface{}, 0)
		for _, msg := range msgs {
			dutyMsgs = append(dutyMsgs, msg)
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

// Return the input for the committee for a given slot of a sequence of duties
func CommitteeInputForSlotInSequencedDuties(numAttestingValidators int, numSyncCommitteeValidators int, numSequencedDuties int, slotIncrement int, currentSlot int, ksMap map[phase0.ValidatorIndex]*TestKeySet, addPostConsensus bool, diffValidatorsForDuties bool) (*types.CommitteeDuty, []*types.SignedSSVMessage) {

	// Validators assigned to the duties
	attValidatorsForDuty, scValidatorsForDuty := selectValidatorsForDuties(numAttestingValidators, numSyncCommitteeValidators, numSequencedDuties, slotIncrement, diffValidatorsForDuties)

	// Duty
	duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsForDuty, scValidatorsForDuty)

	// QBFT and Post-Consensus
	msgs := CommitteeInputForDuty(duty, phase0.Slot(currentSlot), ksMap, addPostConsensus)

	return duty, msgs
}

func CommitteeInputForDuty(duty *types.CommitteeDuty, slot phase0.Slot, ksMap map[phase0.ValidatorIndex]*TestKeySet, addPostConsensus bool) []*types.SignedSSVMessage {

	var sampleKeySet *TestKeySet
	for _, ks := range ksMap {
		sampleKeySet = ks
		break
	}
	msgID := CommitteeMsgID(sampleKeySet)

	ret := make([]*types.SignedSSVMessage, 0)

	// QBFT
	qbftMsgs := SSVDecidingMsgsForHeightWithRoot(sha256.Sum256(TestBeaconVoteByts), TestBeaconVoteByts, msgID, qbft.Height(slot), sampleKeySet)
	ret = append(ret, qbftMsgs...)

	// Post-consensus
	if addPostConsensus {
		for opID := uint64(1); opID <= sampleKeySet.Threshold; opID++ {
			postConsensusMsg := SignPartialSigSSVMessage(sampleKeySet, SSVMsgCommittee(sampleKeySet, nil, PostConsensusCommitteeMsgForDuty(duty, ksMap, opID)))
			ret = append(ret, postConsensusMsg)
		}
	}

	return ret
}

// Committee output

func CommitteeOutputMessagesForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []*types.PartialSignatureMessages {
	return CommitteeOutputMessagesForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, false)
}

func CommitteeOutputMessagesForDutiesWithDifferentValidators(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []*types.PartialSignatureMessages {
	return CommitteeOutputMessagesForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, true)
}

// Returns a list of partial sig output messages for each duty slot
func CommitteeOutputMessagesForDutiesWithParams(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, startingSlot int, diffValidatorsForDuties bool) []*types.PartialSignatureMessages {

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	ret := make([]*types.PartialSignatureMessages, 0)

	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		currentSlot := startingSlot + slotIncrement

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
	return CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, false)
}

func CommitteeBeaconBroadcastedRootsForDutiesWithStartingSlot(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, slot int) []string {
	return CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, slot, false)
}

func CommitteeBeaconBroadcastedRootsForDutiesWithDifferentValidators(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []string {
	return CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties, numAttestingValidators, numSyncCommitteeValidators, TestingDutySlot, true)
}

// Returns a list of beacon roots for committees duties
func CommitteeBeaconBroadcastedRootsForDutiesWithParams(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, startingSlot int, diffValidatorsForDuties bool) []string {

	// KeySet map
	maxValidators := int(math.Max(float64(numAttestingValidators), float64(numSyncCommitteeValidators)))
	ksMap := KeySetMapForValidators(maxValidators)

	ret := make([]string, 0)
	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		currentSlot := startingSlot + slotIncrement

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
