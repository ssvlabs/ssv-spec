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

var SSVMsgAttester = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee))
}

var SSVMsgCommittee = func(ks *TestKeySet, qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	msgIDBytes := CommitteeMsgID(ks)
	var msgID types.MessageID
	copy(msgID[:], msgIDBytes)
	return ssvMsg(qbftMsg, partialSigMsg, msgID)
}

var SSVMsgSyncCommittee = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee))
}

// ==================================================
// Post Consensus
// ==================================================

var PostConsensusCommitteeMsgForDuty = func(duty *types.CommitteeDuty, keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages

	for _, validatorDuty := range duty.ValidatorDuties {

		ks := keySetMap[validatorDuty.ValidatorIndex]

		if validatorDuty.Type == types.BNRoleAttester {
			attData := TestingAttestationDataForValidatorDuty(validatorDuty)
			pSigMsgs := postConsensusAttestationMsgForAttestationData(ks.Shares[id], id, duty.Slot, attData, validatorDuty.ValidatorIndex)
			if ret == nil {
				ret = pSigMsgs
			} else {
				ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
			}
		} else if validatorDuty.Type == types.BNRoleSyncCommittee {
			pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, duty.Slot, false, false, validatorDuty.ValidatorIndex)
			if ret == nil {
				ret = pSigMsgs
			} else {
				ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
			}
		} else {
			panic(fmt.Sprintf("type %v not expected", validatorDuty.Type))
		}
	}

	return ret
}

// ==================================================
// Post Consensus - Attestation
// ==================================================

var PostConsensusAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusAttestationMsgForKeySetWithSlot(keySetMap, id, TestingDutySlotV(version))
}

var PostConsensusAttestationMsgForKeySetWithSlot = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for attestations for each validator in shares
	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		pSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, slot, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	return ret
}

var PostConsensusPartiallyWrongBeaconSigAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationMsgForKeySet(keySetMap, id, version, false, true)
}

var PostConsensusPartiallyWrongRootAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationMsgForKeySet(keySetMap, id, version, true, false)
}

var PostConsensusPartiallyWrongAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		pSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, TestingDutySlotV(version), wrongRootV, wrongBeaconSigV, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}

		msgIndex++
	}
	return ret
}

var PostConsensusWrongAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAttestationMsg(sk, id, TestingDutySlotV(version), true, false, TestingValidatorIndex)
}

var PostConsensusWrongValidatorIndexAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusAttestationMsg(sk, id, TestingDutySlotV(version), true, false, TestingValidatorIndex)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusAttestationMsgForValidatorsIndex = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, validatorIndexes []phase0.ValidatorIndex) *types.PartialSignatureMessages {
	var msg *types.PartialSignatureMessages
	for _, valIdx := range validatorIndexes {
		valIdxMsg := postConsensusAttestationMsg(sk, id, TestingDutySlotV(version), false, false, valIdx)
		if msg == nil {
			msg = valIdxMsg
		} else {
			msg.Messages = append(msg.Messages, valIdxMsg.Messages...)
		}
	}
	return msg
}

var PostConsensusWrongSigAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAttestationMsg(sk, id, TestingDutySlotV(version), false, true, TestingValidatorIndex)
}

var PostConsensusAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAttestationMsg(sk, id, TestingDutySlotV(version), false, false, TestingValidatorIndex)
}

var PostConsensusAttestationTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := postConsensusAttestationMsg(sk, id, TestingDutySlotV(version), false, false, TestingValidatorIndex)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusAttestationTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var postConsensusAttestationMsgForAttestationData = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	slot phase0.Slot,
	attData *phase0.AttestationData,
	validatorIndex phase0.ValidatorIndex,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(attData.Target.Epoch, types.DomainAttester)
	signed, root, _ := signer.SignBeaconObject(attData, d, sk.GetPublicKey().Serialize(), types.DomainAttester)
	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: slot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed,
				SigningRoot:      root,
				Signer:           id,
				ValidatorIndex:   validatorIndex,
			},
		},
	}
	return &msgs
}

var postConsensusAttestationMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	slot phase0.Slot,
	wrongRoot bool,
	wrongBeaconSig bool,
	validatorIndex phase0.ValidatorIndex,
) *types.PartialSignatureMessages {

	sampleAttData := TestingAttestationData(spec.DataVersionPhase0)

	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(sampleAttData.Target.Epoch, types.DomainAttester)

	attData := &phase0.AttestationData{
		Slot:            slot,
		Index:           TestingCommitteeIndex,
		BeaconBlockRoot: sampleAttData.BeaconBlockRoot,
		Source:          sampleAttData.Source,
		Target:          sampleAttData.Target,
	}

	version := VersionBySlot(slot)
	if version >= spec.DataVersionElectra {
		attData.Index = 0
	}

	if wrongRoot {
		attData = TestingWrongAttestationData(spec.DataVersionPhase0)
	}

	signed, root, _ := signer.SignBeaconObject(attData, d, sk.GetPublicKey().Serialize(), types.DomainAttester)

	if wrongBeaconSig {
		signed, _, _ = signer.SignBeaconObject(attData, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainAttester)
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: slot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed,
				SigningRoot:      root,
				Signer:           id,
				ValidatorIndex:   validatorIndex,
			},
		},
	}
	return &msgs
}

// ==================================================
// Post Consensus - Attestation and Sync Committee
// ==================================================

var PostConsensusAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusAttestationAndSyncCommitteeMsgForKeySetWithSlot(keySetMap, id, TestingDutySlotV(version))
}

var PostConsensusAttestationAndSyncCommitteeMsgForKeySetWithSlot = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for attestations for each validator in shares
	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		pSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, slot, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	// Get post consensus for sync committees for each validator in shares
	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, slot, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	return ret
}

var PostConsensusPartiallyWrongBeaconSigAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationAndSyncCommitteeMsgForKeySet(keySetMap, id, version, false, true)
}

var PostConsensusPartiallyWrongRootAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationAndSyncCommitteeMsgForKeySet(keySetMap, id, version, true, false)
}

var PostConsensusPartiallyWrongAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		attPSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, TestingDutySlotV(version), wrongRootV, wrongBeaconSigV, valIdx)
		if ret == nil {
			ret = attPSigMsgs
		} else {
			ret.Messages = append(ret.Messages, attPSigMsgs.Messages...)
		}

		scPSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, TestingDutySlotV(version), wrongRootV, wrongBeaconSigV, valIdx)
		ret.Messages = append(ret.Messages, scPSigMsgs.Messages...)

		msgIndex++
	}
	return ret
}

var PostConsensusAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAttestationAndSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, false)
}

var PostConsensusWrongAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAttestationAndSyncCommitteeMsg(sk, id, TestingDutySlotV(version), true, false)
}

var PostConsensusWrongValidatorIndexAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusAttestationAndSyncCommitteeMsg(sk, id, TestingDutySlotV(version), true, false)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, validatorIndexes []phase0.ValidatorIndex) *types.PartialSignatureMessages {
	var msg *types.PartialSignatureMessages
	for _, valIdx := range validatorIndexes {
		valIdxMsg := postConsensusAttestationAndSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, false)
		for _, m := range valIdxMsg.Messages {
			m.ValidatorIndex = valIdx
		}

		if msg == nil {
			msg = valIdxMsg
		} else {
			msg.Messages = append(msg.Messages, valIdxMsg.Messages...)
		}
	}
	return msg
}

var PostConsensusWrongSigAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAttestationAndSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, true)
}

var PostConsensusAttestationAndSyncCommitteeMsgTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := postConsensusAttestationAndSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, false)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusAttestationAndSyncCommitteeMsgTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var postConsensusAttestationAndSyncCommitteeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	slot phase0.Slot,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	attestationPSigMsg := postConsensusAttestationMsg(sk, id, slot, wrongRoot, wrongBeaconSig, TestingValidatorIndex)
	syncCommitteePSigMessage := postConsensusSyncCommitteeMsg(sk, id, slot, wrongRoot, wrongBeaconSig, TestingValidatorIndex)

	attestationPSigMsg.Messages = append(attestationPSigMsg.Messages, syncCommitteePSigMessage.Messages...)

	return attestationPSigMsg
}

var PreConsensusFailedMsg = func(msgSigner *bls.SecretKey, msgSignerID types.OperatorID) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainRandao)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(TestingDutyEpoch), d, msgSigner.GetPublicKey().Serialize(), types.DomainRandao)

	msg := types.PartialSignatureMessages{
		Type: types.RandaoPartialSig,
		Slot: TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed[:],
				SigningRoot:      root,
				Signer:           msgSignerID,
				ValidatorIndex:   TestingValidatorIndex,
			},
		},
	}
	return &msg
}

// ==================================================
// Post Consensus - Sync Committee
// ==================================================

var PostConsensusSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusSyncCommitteeMsgForKeySetWithSlot(keySetMap, id, TestingDutySlotV(version))
}

var PostConsensusSyncCommitteeMsgForKeySetWithSlot = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for sync committees for each validator in shares
	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, slot, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	return ret
}

var PostConsensusPartiallyWrongBeaconSigSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongSyncCommitteeMsgForKeySet(keySetMap, id, false, true, version)
}

var PostConsensusPartiallyWrongRootSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongSyncCommitteeMsgForKeySet(keySetMap, id, true, false, version)
}

var PostConsensusPartiallyWrongSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, wrongRoot bool, wrongBeaconSig bool, version spec.DataVersion) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	for _, valKs := range SortedMapKeys(keySetMap) {
		ks := valKs.Value
		valIdx := valKs.Key

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, TestingDutySlotV(version), wrongRootV, wrongBeaconSigV, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}

		msgIndex++
	}
	return ret
}

var PostConsensusSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, false, TestingValidatorIndex)
}

var PostConsensusSyncCommitteeTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := postConsensusSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, false, TestingValidatorIndex)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusSyncCommitteeTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var PostConsensusWrongSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeMsg(sk, id, TestingDutySlotV(version), true, false, TestingValidatorIndex)
}

var PostConsensusWrongValidatorIndexSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusSyncCommitteeMsg(sk, id, TestingDutySlotV(version), true, false, TestingValidatorIndex)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusSyncCommitteeMsgForValidatorsIndex = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion, validatorIndexes []phase0.ValidatorIndex) *types.PartialSignatureMessages {
	var msg *types.PartialSignatureMessages
	for _, valIdx := range validatorIndexes {
		valIdxMsg := postConsensusSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, false, valIdx)
		if msg == nil {
			msg = valIdxMsg
		} else {
			msg.Messages = append(msg.Messages, valIdxMsg.Messages...)
		}
	}
	return msg
}

var PostConsensusWrongSigSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeMsg(sk, id, TestingDutySlotV(version), false, true, TestingValidatorIndex)
}

var postConsensusSyncCommitteeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	slot phase0.Slot,
	wrongRoot bool,
	wrongBeaconSig bool,
	validatorIndex phase0.ValidatorIndex,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(phase0.Epoch(slot/32), types.DomainSyncCommittee)
	blockRoot := TestingSyncCommitteeBlockRoot
	if wrongRoot {
		blockRoot = TestingSyncCommitteeWrongBlockRoot
	}
	signed, root, _ := signer.SignBeaconObject(types.SSZBytes(blockRoot[:]), d, sk.GetPublicKey().Serialize(), types.DomainSyncCommittee)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(types.SSZBytes(blockRoot[:]), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainSyncCommittee)
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: slot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed,
				SigningRoot:      root,
				Signer:           id,
				ValidatorIndex:   validatorIndex,
			},
		},
	}
	return &msgs
}
