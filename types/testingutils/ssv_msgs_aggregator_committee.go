package testingutils

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

func AggregatorCommitteeMsgIDForKeySet(ks *TestKeySet) types.MessageID {
	msgIDBytes := AggregatorCommitteeMsgID(ks)
	var msgID types.MessageID
	copy(msgID[:], msgIDBytes)
	return msgID
}

var TestingAggregatorCommitteeMsgID = func() types.MessageID {
	return AggregatorCommitteeMsgIDForKeySet(Testing4SharesSet())
}()

var SSVMsgAggregatorCommittee = func(keySet *TestKeySet, qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	// For aggregator committee, we use committee ID since it manages multiple validators
	return ssvMsg(qbftMsg, partialSigMsg, AggregatorCommitteeMsgIDForKeySet(keySet))
}

// ==================================================
// Pre-Consensus Messages for AggregatorCommittee
// ==================================================

// PreConsensusAggregatorCommitteeMsg creates a pre-consensus message for aggregator committee (selection proofs)
// This is used for single validator test cases with aggregator role
var PreConsensusAggregatorCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// The aggregator committee runner broadcasts AggregatorCommitteePartialSig type
	// containing selection proofs for aggregator duties
	selectionProofMsg := PreConsensusSelectionProofMsg(sk, sk, id, id, version)

	// Convert to AggregatorCommitteePartialSig type
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: selectionProofMsg.Messages,
	}
}

// ==================================================
// Post-Consensus Messages for AggregatorCommittee
// ==================================================

// PostConsensusAggregatorCommitteeMsg creates a post-consensus message for aggregator committee
// This is used for single validator test cases with aggregator role
var PostConsensusAggregatorCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// For aggregator committee post-consensus, we use the same message as regular aggregator
	return PostConsensusAggregatorMsg(sk, id, version)
}

// PreConsensusAggregatorCommitteeSyncCommitteeMsg creates a pre-consensus message for sync committee contribution
// This is used for single validator test cases with sync committee role
var PreConsensusAggregatorCommitteeSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	// The aggregator committee runner broadcasts AggregatorCommitteePartialSig type
	// containing contribution proofs for sync committee duties
	contribProofMsg := PreConsensusContributionProofMsg(sk, sk, id, id)

	// Convert to AggregatorCommitteePartialSig type
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     contribProofMsg.Slot,
		Messages: contribProofMsg.Messages,
	}
}

// PostConsensusAggregatorCommitteeSyncCommitteeMsg creates a post-consensus message for sync committee
// This is used for single validator test cases with sync committee role
var PostConsensusAggregatorCommitteeSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	// For sync committee post-consensus, we use the sync committee contribution message
	return PostConsensusSyncCommitteeContributionMsg(sk, id, keySet)
}

// PreConsensusAggregatorCommitteeMixedMsg creates pre-consensus messages for mixed duties (aggregator + sync committee)
// The aggregator committee runner broadcasts both selection proofs and contribution proofs in one message
var PreConsensusAggregatorCommitteeMixedMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// Get aggregator selection proof
	selectionProofMsg := PreConsensusSelectionProofMsg(sk, sk, id, id, version)

	// Get sync committee contribution proofs
	contribProofMsg := PreConsensusContributionProofMsgWithSlot(sk, sk, id, id, TestingDutySlotV(version))

	// Combine into a single AggregatorCommitteePartialSig message
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: append(selectionProofMsg.Messages, contribProofMsg.Messages...),
	}
}

var PreConsensusAggregatorCommitteeMixedMsgTooManyRoots = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// Get aggregator selection proof
	selectionProofMsg := PreConsensusSelectionProofTooManyRootsMsg(sk, sk, id, id, version)

	// Get sync committee contribution proofs
	contribProofMsg := PreConsensusContributionProofTooManyRootsMsg(sk, sk, id, id, TestingDutySlotV(version))

	// Combine into a single AggregatorCommitteePartialSig message
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: append(selectionProofMsg.Messages, contribProofMsg.Messages...),
	}
}
var PreConsensusAggregatorCommitteeMixedMsgTooFewRoots = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// Get aggregator selection proof
	selectionProofMsg := PreConsensusSelectionProofTooFewRootsMsg(sk, sk, id, id, version)

	// Get sync committee contribution proofs
	contribProofMsg := PreConsensusContributionProofTooFewRootsMsg(sk, sk, id, id, TestingDutySlotV(version))

	// Combine into a single AggregatorCommitteePartialSig message
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: append(selectionProofMsg.Messages, contribProofMsg.Messages...),
	}
}

var PreConsensusAggregatorCommitteeMixedMsgInconsistentBeaconSigners = func(sk *bls.SecretKey, id, id2 types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// Get aggregator selection proof
	selectionProofMsg := PreConsensusSelectionProofMsg(sk, sk, id, id, version)

	// Get sync committee contribution proofs
	contribProofMsg := PreConsensusContributionProofMsgInconsistentBeaconSigners(sk, sk, id, id, id2, TestingDutySlotV(version))

	// Combine into a single AggregatorCommitteePartialSig message
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: append(selectionProofMsg.Messages, contribProofMsg.Messages...),
	}
}

var PreConsensusAggregatorCommitteeMixedMsgWrongRoot = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// Get aggregator selection proof
	selectionProofMsg := PreConsensusSelectionProofWrongRootSigMsg(sk, sk, id, id, version)

	// Get sync committee contribution proofs
	contribProofMsg := PreConsensusContributionProofWrongBeaconRootMsg(sk, sk, id, id, TestingDutySlotV(version))

	// Combine into a single AggregatorCommitteePartialSig message
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: append(selectionProofMsg.Messages, contribProofMsg.Messages...),
	}
}

var PreConsensusAggregatorCommitteeMixedMsgWrongBeaconSig = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	// Get aggregator selection proof
	selectionProofMsg := PreConsensusSelectionProofWrongBeaconSigMsg(sk, sk, id, id, version)

	// Get sync committee contribution proofs
	contribProofMsg := PreConsensusContributionProofWrongBeaconSigMsg(sk, sk, id, id, TestingDutySlotV(version))

	// Combine into a single AggregatorCommitteePartialSig message
	return &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     selectionProofMsg.Slot,
		Messages: append(selectionProofMsg.Messages, contribProofMsg.Messages...),
	}
}

// PostConsensusAggregatorCommitteeMixedMsg creates post-consensus messages for mixed duties
// This combines signatures for both aggregator and sync committee roles
var PostConsensusAggregatorCommitteeMixedMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusAggregatorMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusSyncCommitteeContributionMsgWithSlot(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

var PostConsensusAggregatorCommitteeMixedMsgWrongValIdx = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusWrongValidatorIndexAggregatorMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusWrongValidatorIndexSyncCommitteeContributionMsg(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

var PostConsensusAggregatorCommitteeMixedWrongSigMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusWrongSigAggregatorMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusWrongSigSyncCommitteeContributionMsg(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

var PostConsensusAggregatorCommitteeMixedWrongMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusWrongAggregatorMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusWrongSyncCommitteeContributionMsg(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

var PostConsensusAggregatorCommitteeMixedMsgTooFewRoots = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusAggregatorTooFewRootsMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusSyncCommitteeContributionTooFewRootsMsg(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

var PostConsensusAggregatorCommitteeMixedMsgTooManyRoots = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusAggregatorTooManyRootsMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusSyncCommitteeContributionTooManyRootsMsg(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

var PostConsensusAggregatorCommitteeWrongOrderMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, keySet *TestKeySet) *types.PartialSignatureMessages {
	// Get aggregator post-consensus message
	aggMsg := PostConsensusAggregatorMsg(sk, id, version)

	// Get sync committee post-consensus messages
	syncMsg := PostConsensusSyncCommitteeContributionWrongOrderMsg(sk, id, keySet, TestingDutySlotV(version))

	// Combine messages - all with PostConsensusPartialSig type
	aggMsg.Messages = append(aggMsg.Messages, syncMsg.Messages...)

	return aggMsg
}

// PostConsensusAggregatorCommitteeMsgForDuty creates post-consensus messages for all validators in the duty
var PostConsensusAggregatorCommitteeMsgForDuty = func(duty *types.AggregatorCommitteeDuty, keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	var ret *types.PartialSignatureMessages

	for _, validatorDuty := range duty.ValidatorDuties {
		ks := keySetMap[validatorDuty.ValidatorIndex]
		if ks == nil {
			continue
		}

		// Only include messages for validators where this operator has a share
		if ks.Shares[id] == nil {
			continue
		}

		switch validatorDuty.Type {
		case types.BNRoleAggregator:
			pSigMsgs := PostConsensusAggregatorMsgWithValidatorIndex(ks.Shares[id], id, validatorDuty.ValidatorIndex, version)
			if ret == nil {
				ret = pSigMsgs
			} else {
				ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
			}
		case types.BNRoleSyncCommitteeContribution:
			pSigMsgs := PostConsensusCustomSlotSyncCommitteeContributionMsgWithValidatorIndex(ks.Shares[id], validatorDuty.ValidatorIndex, id, ks, validatorDuty.Slot)
			if ret == nil {
				ret = pSigMsgs
			} else {
				ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
			}
		default:
			panic(fmt.Sprintf("type %v not expected", validatorDuty.Type))
		}
	}

	return ret
}

// PreConsensusAggregatorCommitteeMsgForDuty creates pre-consensus messages for all validators in the duty
// This handles selection proofs for both aggregator and sync committee roles
var PreConsensusAggregatorCommitteeMsgForDuty = func(duty *types.AggregatorCommitteeDuty, keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	var ret *types.PartialSignatureMessages

	for _, validatorDuty := range duty.ValidatorDuties {
		ks := keySetMap[validatorDuty.ValidatorIndex]
		if ks == nil {
			continue
		}

		switch validatorDuty.Type {
		case types.BNRoleAggregator:
			// Add aggregator selection proof with SelectionProofPartialSig type
			selectionProofMsg := PreConsensusSelectionProofMsgWithValidatorIndex(ks.Shares[id], ks.Shares[id], validatorDuty.ValidatorIndex, id, id, version)
			if ret == nil {
				ret = selectionProofMsg
			} else {
				ret.Messages = append(ret.Messages, selectionProofMsg.Messages...)
			}

		case types.BNRoleSyncCommitteeContribution:
			// Add sync committee contribution proofs with ContributionProofsPartialSig type
			contributionProofMsg := PreConsensusCustomSlotContributionProofMsgWithValidatorIndex(ks.Shares[id], ks.Shares[id], validatorDuty.ValidatorIndex, id, id, duty.Slot)
			if ret == nil {
				ret = contributionProofMsg
			} else {
				ret.Messages = append(ret.Messages, contributionProofMsg.Messages...)
			}
		default:
			panic(fmt.Sprintf("type %v not expected", validatorDuty.Type))
		}
	}

	// Convert to AggregatorCommitteePartialSig type as used by the runner
	if ret != nil {
		ret.Type = types.AggregatorCommitteePartialSig
	}

	return ret
}

var PostConsensusAggMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusAggCommitteeMsgForKeySet(keySetMap, id, version, true, false)
}
var PostConsensusSCCMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusAggCommitteeMsgForKeySet(keySetMap, id, version, false, true)
}
var PostConsensusAggAndSCCMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusAggCommitteeMsgForKeySet(keySetMap, id, version, true, true)
}

var PostConsensusAggCommitteeMsgForKeySet = func(
	keySetMap map[phase0.ValidatorIndex]*TestKeySet,
	id types.OperatorID,
	version spec.DataVersion,
	includeAgg bool,
	includeSCC bool,
) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for attestations for each validator in shares
	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key

		if includeAgg {
			aggMsg := postConsensusAggregatorMsg(ks.Shares[id], id, valIdx, false, false, version)
			if ret == nil {
				ret = aggMsg
			} else {
				ret.Messages = append(ret.Messages, aggMsg.Messages...)
			}
		}

		if includeSCC {
			sccMsg := postConsensusSyncCommitteeContributionMsg(ks.Shares[id], id, valIdx, TestingDutySlotV(version), ks, false, false, false)
			if ret == nil {
				ret = sccMsg
			} else {
				ret.Messages = append(ret.Messages, sccMsg.Messages...)
			}
		}
	}
	return ret
}

var PostConsensusPartiallyWrongAggMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAggCommmitteeMsgForKeySet(keySetMap, id, version, wrongRoot, wrongBeaconSig, true, false)
}
var PostConsensusPartiallyWrongSCCMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAggCommmitteeMsgForKeySet(keySetMap, id, version, wrongRoot, wrongBeaconSig, false, true)
}
var PostConsensusPartiallyWrongAggAndSCCMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAggCommmitteeMsgForKeySet(keySetMap, id, version, wrongRoot, wrongBeaconSig, true, true)
}

var PostConsensusPartiallyWrongAggCommmitteeMsgForKeySet = func(
	keySetMap map[phase0.ValidatorIndex]*TestKeySet,
	id types.OperatorID,
	version spec.DataVersion,
	wrongRoot bool,
	wrongBeaconSig bool,
	includeAgg bool,
	includeSCC bool,
) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		if includeAgg {

			attPSigMsgs := postConsensusAggregatorMsg(ks.Shares[id], id, valIdx, wrongRootV, wrongBeaconSigV, version)
			if ret == nil {
				ret = attPSigMsgs
			} else {
				ret.Messages = append(ret.Messages, attPSigMsgs.Messages...)
			}
		}

		if includeSCC {
			scPSigMsgs := postConsensusSyncCommitteeContributionMsg(ks.Shares[id], id, valIdx, TestingDutySlotV(version), ks, wrongRootV, wrongBeaconSigV, false)
			if ret == nil {
				ret = scPSigMsgs
			} else {
				ret.Messages = append(ret.Messages, scPSigMsgs.Messages...)
			}
		}

		msgIndex++
	}
	return ret
}

// PreConsensusAggregatorCommitteeMsgForDutySorted creates pre-consensus messages for all validators in the duty
// and sorts them by (validator index, signing root)
func PreConsensusAggregatorCommitteeMsgForDutySorted(duty *types.AggregatorCommitteeDuty, keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, id, version)

	if ret == nil || len(ret.Messages) == 0 {
		return ret
	}
	ret.Sort()
	return ret
}

// PostConsensusAggregatorCommitteeMsgForDutySorted creates post-consensus messages for all validators in the duty
// and sorts them by (validator index, signing root)
func PostConsensusAggregatorCommitteeMsgForDutySorted(duty *types.AggregatorCommitteeDuty, keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := PostConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, id, version)
	if ret == nil || len(ret.Messages) == 0 {
		return ret
	}
	ret.Sort()
	return ret
}
