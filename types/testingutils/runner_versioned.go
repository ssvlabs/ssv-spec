package testingutils

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

var SSVDecidingMsgsForCommitteeRunner = func(beaconVote *types.BeaconVote, ks *TestKeySet, height qbft.Height) []*types.SignedSSVMessage {
	id := CommitteeMsgID(ks)

	// consensus
	qbftMsgs := SSVDecidingMsgsForHeightAndBeaconVote(beaconVote, id[:], height, ks)
	return qbftMsgs
}

var SSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.RunnerRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	signedF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, _ := partialSigMsg.Encode()
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewTestingOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
		if err != nil {
			panic(err)
		}
		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{signer},
			Signatures:  [][]byte{sig},
			SSVMessage:  ssvMsg,
		}
	}

	// pre consensus msgs
	base := make([]*types.SignedSSVMessage, 0)
	if role == types.RoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedF(PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
		}
	}
	if role == types.RoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedF(PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.RoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedF(PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)
	return base
}

var ExpectedSSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.RunnerRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, _ := partialSigMsg.Encode()

		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewTestingOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
		if err != nil {
			panic(err)
		}
		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{signer},
			Signatures:  [][]byte{sig},
			SSVMessage:  ssvMsg,
		}
	}

	// pre consensus msgs
	base := make([]*types.SignedSSVMessage, 0)
	if role == types.RoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
		}
	}
	if role == types.RoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.RoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVExpectedDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)
	return base
}

// Returns a list of [Duty, qbft messages...] for each duty
func CommitteeInputForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int, addPostConsensus bool, shuffle bool) []interface{} {

	// Attestation validators and KeySet map
	attValidatorsIndexList := ValidatorIndexList(numAttestingValidators)
	attKsMap := KeySetMapForValidators(numAttestingValidators)

	// Sync committee validators and KeySet map
	scValidatorsIndexList := ValidatorIndexList(numSyncCommitteeValidators)
	scKsMap := KeySetMapForValidators(numSyncCommitteeValidators)

	// Joint KeySet map
	jointMap := attKsMap
	for valIdx, valKS := range scKsMap {
		jointMap[valIdx] = valKS
	}

	// Sample KeySet
	var sampleKeySet *TestKeySet
	for _, ks := range jointMap {
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

		// Duty
		duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList)
		dutiesCommands = append(dutiesCommands, duty)

		// QBFT
		for _, msg := range SSVDecidingMsgsForHeightWithRoot(sha256.Sum256(TestBeaconVoteByts), TestBeaconVoteByts, msgID, qbft.Height(currentSlot), sampleKeySet) {
			dutyMsgs = append(dutyMsgs, msg)
		}

		// Post-consensus
		if addPostConsensus {
			for opID := uint64(1); opID <= sampleKeySet.Threshold; opID++ {
				postConsensusMsg := SignPartialSigSSVMessage(sampleKeySet, SSVMsgCommittee(sampleKeySet, nil, PostConsensusCommitteeMsgForDuty(duty, jointMap, opID)))
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

// Returns a list of partial sig output messages for each duty slot
func CommitteeOutputMessagesForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []*types.PartialSignatureMessages {

	// Attestation validators and KeySet map
	attValidatorsIndexList := ValidatorIndexList(numAttestingValidators)
	attKsMap := KeySetMapForValidators(numAttestingValidators)

	// Sync committee validators and KeySet map
	scValidatorsIndexList := ValidatorIndexList(numSyncCommitteeValidators)
	scKsMap := KeySetMapForValidators(numSyncCommitteeValidators)

	// Joint KeySet map
	jointMap := attKsMap
	for valIdx, valKS := range scKsMap {
		jointMap[valIdx] = valKS
	}

	ret := make([]*types.PartialSignatureMessages, 0)

	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		currentSlot := TestingDutySlot + slotIncrement

		duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList)
		postConsensusMsg := PostConsensusCommitteeMsgForDuty(duty, jointMap, 1)

		// Add post consensus
		ret = append(ret, postConsensusMsg)
	}

	return ret
}

// Returns a list of beacon roots for committees duties
func CommitteeBeaconBroadcastedRootsForDuties(numSequencedDuties int, numAttestingValidators int, numSyncCommitteeValidators int) []string {

	// Attestation validators and KeySet map
	attValidatorsIndexList := ValidatorIndexList(numAttestingValidators)
	attKsMap := KeySetMapForValidators(numAttestingValidators)

	// Sync committee validators and KeySet map
	scValidatorsIndexList := ValidatorIndexList(numSyncCommitteeValidators)
	scKsMap := KeySetMapForValidators(numSyncCommitteeValidators)

	// Joint KeySet map
	jointMap := attKsMap
	for valIdx, valKS := range scKsMap {
		jointMap[valIdx] = valKS
	}

	ret := make([]string, 0)
	for slotIncrement := 0; slotIncrement < numSequencedDuties; slotIncrement++ {

		currentSlot := TestingDutySlot + slotIncrement

		duty := TestingCommitteeDuty(phase0.Slot(currentSlot), attValidatorsIndexList, scValidatorsIndexList)

		beaconRoots := TestingSignedCommitteeBeaconObjectSSZRoot(duty, jointMap)

		// Add beacon roots
		ret = append(ret, beaconRoots...)
	}
	return ret
}

// Returns a list of beacon roots for committees duties
func CommitteeBeaconBroadcastedRootsForDuty(slot phase0.Slot, numAttestingValidators int, numSyncCommitteeValidators int) []string {

	// Attestation validators and KeySet map
	attValidatorsIndexList := ValidatorIndexList(numAttestingValidators)
	attKsMap := KeySetMapForValidators(numAttestingValidators)

	// Sync committee validators and KeySet map
	scValidatorsIndexList := ValidatorIndexList(numSyncCommitteeValidators)
	scKsMap := KeySetMapForValidators(numSyncCommitteeValidators)

	// Joint KeySet map
	jointMap := attKsMap
	for valIdx, valKS := range scKsMap {
		jointMap[valIdx] = valKS
	}

	ret := make([]string, 0)

	duty := TestingCommitteeDuty(slot, attValidatorsIndexList, scValidatorsIndexList)

	beaconRoots := TestingSignedCommitteeBeaconObjectSSZRoot(duty, jointMap)

	// Add beacon roots
	ret = append(ret, beaconRoots...)

	return ret
}
