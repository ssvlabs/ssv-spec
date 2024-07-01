package testingutils

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingSSVDomainType = types.JatoTestnet
var TestingForkData = types.ForkData{Epoch: TestingDutyEpoch, Domain: TestingSSVDomainType}
var CommitteeMsgID = func(keySet *TestKeySet) []byte {

	// Identifier
	committee := make([]uint64, 0)
	for _, op := range keySet.Committee() {
		committee = append(committee, op.Signer)
	}
	committeeID := types.GetCommitteeID(committee)

	ret := types.NewMsgID(TestingSSVDomainType, committeeID[:], types.RoleCommittee)
	return ret[:]
}
var AttesterMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee)
	return ret[:]
}()
var ProposerMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleProposer)
	return ret[:]
}()
var AggregatorMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleAggregator)
	return ret[:]
}()
var SyncCommitteeMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee)
	return ret[:]
}()
var SyncCommitteeContributionMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleSyncCommitteeContribution)
	return ret[:]
}()
var ValidatorRegistrationMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleValidatorRegistration)
	return ret[:]
}()
var VoluntaryExitMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleVoluntaryExit)
	return ret[:]
}()

var EncodeConsensusDataTest = func(cd *types.ValidatorConsensusData) []byte {
	encodedCD, _ := cd.Encode()
	return encodedCD
}

// BeaconVote data - Committee Runner

var TestBeaconVoteByts, _ = TestBeaconVote.Encode()

var TestBeaconVoteNextEpochByts, _ = TestBeaconVoteNextEpoch.Encode()

var TestWrongBeaconVoteByts, _ = TestWrongBeaconVote.Encode()

// ValidatorConsensusData - Attester

var TestAttesterConsensusData = &types.ValidatorConsensusData{
	Duty:    *TestingAttesterDuty.ValidatorDuties[0],
	DataSSZ: TestingAttestationDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

var TestAttesterNextEpochConsensusData = &types.ValidatorConsensusData{
	Duty:    *TestingAttesterDutyNextEpoch.ValidatorDuties[0],
	DataSSZ: TestingAttestationNextEpochDataBytes,
	Version: spec.DataVersionPhase0,
}

var TestingAttesterNextEpochConsensusDataByts, _ = TestAttesterNextEpochConsensusData.Encode()

// ValidatorConsensusData - Aggregator

var TestAggregatorConsensusData = &types.ValidatorConsensusData{
	Duty:    TestingAggregatorDuty,
	DataSSZ: TestingAggregateAndProofBytes,
	Version: spec.DataVersionPhase0,
}
var TestAggregatorConsensusDataByts, _ = TestAggregatorConsensusData.Encode()

// ValidatorConsensusData - Sync Committee

var TestSyncCommitteeConsensusData = &types.ValidatorConsensusData{
	Duty:    *TestingSyncCommitteeDuty.ValidatorDuties[0],
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeConsensusDataByts, _ = TestSyncCommitteeConsensusData.Encode()

var TestSyncCommitteeNextEpochConsensusData = &types.ValidatorConsensusData{
	Duty:    *TestingSyncCommitteeDutyNextEpoch.ValidatorDuties[0],
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: spec.DataVersionPhase0,
}

var TestSyncCommitteeNextEpochConsensusDataByts, _ = TestSyncCommitteeNextEpochConsensusData.Encode()

// ValidatorConsensusData - Sync Committee Contribution

var TestSyncCommitteeContributionConsensusData = &types.ValidatorConsensusData{
	Duty:    TestingSyncCommitteeContributionDuty,
	DataSSZ: TestingContributionsDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()
var TestSyncCommitteeContributionConsensusDataRoot = func() [32]byte {
	return sha256.Sum256(TestSyncCommitteeContributionConsensusDataByts)
}()

var TestConsensusUnkownDutyTypeData = &types.ValidatorConsensusData{
	Duty:    TestingUnknownDutyType,
	DataSSZ: TestingAttestationDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var TestConsensusWrongDutyPKData = &types.ValidatorConsensusData{
	Duty:    TestingWrongDutyPK,
	DataSSZ: TestingAttestationDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestConsensusWrongDutyPKDataByts, _ = TestConsensusWrongDutyPKData.Encode()

var SSVMsgAttester = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee))
}

var SSVMsgWrongID = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingWrongValidatorPubKey[:], types.RoleCommittee))
}

var SSVMsgCommittee = func(ks *TestKeySet, qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	msgIDBytes := CommitteeMsgID(ks)
	var msgID types.MessageID
	copy(msgID[:], msgIDBytes)
	return ssvMsg(qbftMsg, partialSigMsg, msgID)
}

var SSVMsgProposer = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleProposer))
}

var SSVMsgAggregator = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleAggregator))
}

var SSVMsgSyncCommittee = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee))
}

var SSVMsgSyncCommitteeContribution = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleSyncCommitteeContribution))
}

var SSVMsgValidatorRegistration = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleValidatorRegistration))
}

var SSVMsgVoluntaryExit = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleVoluntaryExit))
}

var ssvMsg = func(qbftMsg *types.SignedSSVMessage, postMsg *types.PartialSignatureMessages, msgID types.MessageID) *types.SSVMessage {

	if qbftMsg != nil {
		return &types.SSVMessage{
			MsgType: qbftMsg.SSVMessage.MsgType,
			MsgID:   msgID,
			Data:    qbftMsg.SSVMessage.Data,
		}
	}

	if postMsg != nil {
		msgType := types.SSVPartialSignatureMsgType
		data, err := postMsg.Encode()
		if err != nil {
			panic(err)
		}
		return &types.SSVMessage{
			MsgType: msgType,
			MsgID:   msgID,
			Data:    data,
		}
	}

	panic("msg type undefined")
}

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

var PostConsensusAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for attestations for each validator in shares
	for valIdx, ks := range keySetMap {
		pSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, height, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	return ret
}

var PostConsensusPartiallyWrongBeaconSigAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationMsgForKeySet(keySetMap, id, height, false, true)
}

var PostConsensusPartiallyWrongRootAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationMsgForKeySet(keySetMap, id, height, true, false)
}

var PostConsensusPartiallyWrongAttestationMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	validatorIndexes := make([]phase0.ValidatorIndex, 0)
	for valIdx := range keySetMap {
		validatorIndexes = append(validatorIndexes, valIdx)
	}
	sort.Slice(validatorIndexes, func(i, j int) bool {
		return validatorIndexes[i] < validatorIndexes[j]
	})

	for _, valIdx := range validatorIndexes {
		ks, ok := keySetMap[valIdx]
		if !ok {
			panic("validator index not in key set map")
		}

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		pSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, height, wrongRootV, wrongBeaconSigV, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}

		msgIndex++
	}
	return ret
}

var PostConsensusWrongAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return postConsensusAttestationMsg(sk, id, height, true, false, TestingValidatorIndex)
}

var PostConsensusWrongValidatorIndexAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	msg := postConsensusAttestationMsg(sk, id, height, true, false, TestingValidatorIndex)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return postConsensusAttestationMsg(sk, id, height, false, true, TestingValidatorIndex)
}

var PostConsensusAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return postConsensusAttestationMsg(sk, id, height, false, false, TestingValidatorIndex)
}

var PostConsensusAttestationTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	ret := postConsensusAttestationMsg(sk, id, height, false, false, TestingValidatorIndex)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusAttestationTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
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
	height qbft.Height,
	wrongRoot bool,
	wrongBeaconSig bool,
	validatorIndex phase0.ValidatorIndex,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingAttestationData.Target.Epoch, types.DomainAttester)

	attData := &phase0.AttestationData{
		Slot:            phase0.Slot(height),
		Index:           TestingAttestationData.Index,
		BeaconBlockRoot: TestingAttestationData.BeaconBlockRoot,
		Source:          TestingAttestationData.Source,
		Target:          TestingAttestationData.Target,
	}
	if wrongRoot {
		attData = TestingWrongAttestationData
	}

	signed, root, _ := signer.SignBeaconObject(attData, d, sk.GetPublicKey().Serialize(), types.DomainAttester)

	if wrongBeaconSig {
		signed, _, _ = signer.SignBeaconObject(attData, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainAttester)
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: phase0.Slot(height),
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

// Post Consensus - Attestation and Sync Committee

var PostConsensusAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for attestations for each validator in shares
	for valIdx, ks := range keySetMap {
		pSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, height, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	// Get post consensus for sync committees for each validator in shares
	for valIdx, ks := range keySetMap {
		pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, phase0.Slot(height), false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	return ret
}

var PostConsensusPartiallyWrongBeaconSigAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationAndSyncCommitteeMsgForKeySet(keySetMap, id, height, false, true)
}

var PostConsensusPartiallyWrongRootAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongAttestationAndSyncCommitteeMsgForKeySet(keySetMap, id, height, true, false)
}

var PostConsensusPartiallyWrongAttestationAndSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, height qbft.Height, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	validatorIndexes := make([]phase0.ValidatorIndex, 0)
	for valIdx := range keySetMap {
		validatorIndexes = append(validatorIndexes, valIdx)
	}
	sort.Slice(validatorIndexes, func(i, j int) bool {
		return validatorIndexes[i] < validatorIndexes[j]
	})

	for _, valIdx := range validatorIndexes {
		ks, ok := keySetMap[valIdx]
		if !ok {
			panic("validator index not in key set map")
		}

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		attPSigMsgs := postConsensusAttestationMsg(ks.Shares[id], id, height, wrongRootV, wrongBeaconSigV, valIdx)
		if ret == nil {
			ret = attPSigMsgs
		} else {
			ret.Messages = append(ret.Messages, attPSigMsgs.Messages...)
		}

		scPSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, phase0.Slot(height), wrongRootV, wrongBeaconSigV, valIdx)
		ret.Messages = append(ret.Messages, scPSigMsgs.Messages...)

		msgIndex++
	}
	return ret
}

var PostConsensusAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return postConsensusAttestationAndSyncCommitteeMsg(sk, id, height, false, false)
}

var PostConsensusWrongAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return postConsensusAttestationAndSyncCommitteeMsg(sk, id, height, true, false)
}

var PostConsensusWrongValidatorIndexAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	msg := postConsensusAttestationAndSyncCommitteeMsg(sk, id, height, true, false)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigAttestationAndSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	return postConsensusAttestationAndSyncCommitteeMsg(sk, id, height, false, true)
}

var PostConsensusAttestationAndSyncCommitteeMsgTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.PartialSignatureMessages {
	ret := postConsensusAttestationAndSyncCommitteeMsg(sk, id, height, false, false)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusAttestationAndSyncCommitteeMsgTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var postConsensusAttestationAndSyncCommitteeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	height qbft.Height,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	attestationPSigMsg := postConsensusAttestationMsg(sk, id, height, wrongRoot, wrongBeaconSig, TestingValidatorIndex)
	syncCommitteePSigMessage := postConsensusSyncCommitteeMsg(sk, id, phase0.Slot(height), wrongRoot, wrongBeaconSig, TestingValidatorIndex)

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

var PreConsensusRandaoMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 1, false)
}

var randaoMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	epoch phase0.Epoch,
	msgCnt int,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(epoch, types.DomainRandao)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(epoch), d, sk.GetPublicKey().Serialize(), types.DomainRandao)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(types.SSZUint64(TestingDutyEpoch), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainRandao)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
		if wrongRoot {
			msg.SigningRoot = [32]byte{}
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	return &msgs
}

var PreConsensusSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return PreConsensusCustomSlotSelectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot)
}

var PreConsensusSelectionProofWrongBeaconSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, 1, true)
}

var PreConsensusSelectionProofNextEpochMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot2, TestingDutySlot2, 1, false)
}

var PreConsensusSelectionProofTooManyRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, 3, false)
}

var PreConsensusSelectionProofTooFewRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, 0, false)
}

var PreConsensusCustomSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, slot, TestingDutySlot, 1, false)
}

var PreConsensusWrongMsgSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, 1, false)
}

var selectionProofMsg = func(
	sk *bls.SecretKey,
	beaconsk *bls.SecretKey,
	id types.OperatorID,
	beaconid types.OperatorID,
	slot phase0.Slot,
	msgSlot phase0.Slot,
	msgCnt int,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainSelectionProof)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(slot), d, beaconsk.GetPublicKey().Serialize(), types.DomainSelectionProof)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(types.SSZUint64(slot), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainSelectionProof)
	}

	_msgs := make([]*types.PartialSignatureMessage, 0)
	for i := 0; i < msgCnt; i++ {
		_msgs = append(_msgs, &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           beaconid,
			ValidatorIndex:   TestingValidatorIndex,
		})
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.SelectionProofPartialSig,
		Slot:     TestingDutySlot,
		Messages: _msgs,
	}
	return &msgs
}

var PreConsensusValidatorRegistrationMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, true)
}

var PreConsensusValidatorRegistrationWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var validatorRegistrationMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainApplicationBuilder)

	signed, root, _ := signer.SignBeaconObject(TestingValidatorRegistrationBySlot(slot), d,
		beaconSK.GetPublicKey().Serialize(),
		types.DomainApplicationBuilder)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingValidatorRegistrationWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainApplicationBuilder)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingValidatorRegistration, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainApplicationBuilder)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.ValidatorRegistrationPartialSig,
		Slot:     slot,
		Messages: []*types.PartialSignatureMessage{},
	}

	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           beaconID,
			ValidatorIndex:   TestingValidatorIndex,
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	msg := &types.PartialSignatureMessage{
		PartialSignature: signed[:],
		SigningRoot:      root,
		Signer:           id,
		ValidatorIndex:   TestingValidatorIndex,
	}
	if wrongRoot {
		msg.SigningRoot = [32]byte{}
	}

	return &msgs
}

var PreConsensusVoluntaryExitMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var PreConsensusVoluntaryExitTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, true)
}

var PreConsensusVoluntaryExitWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var VoluntaryExitMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainVoluntaryExit)

	signed, root, _ := signer.SignBeaconObject(TestingVoluntaryExitBySlot(slot), d,
		beaconSK.GetPublicKey().Serialize(),
		types.DomainVoluntaryExit)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingVoluntaryExitWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainVoluntaryExit)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingVoluntaryExit, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainVoluntaryExit)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.VoluntaryExitPartialSig,
		Slot:     slot,
		Messages: []*types.PartialSignatureMessage{},
	}

	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           beaconID,
			ValidatorIndex:   TestingValidatorIndex,
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	msg := &types.PartialSignatureMessage{
		PartialSignature: signed[:],
		SigningRoot:      root,
		Signer:           id,
		ValidatorIndex:   TestingValidatorIndex,
	}
	if wrongRoot {
		msg.SigningRoot = [32]byte{}
	}

	return &msgs
}

var PostConsensusAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, false, false)
}

var PostConsensusAggregatorTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	ret := postConsensusAggregatorMsg(sk, id, false, false)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusAggregatorTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var PostConsensusWrongAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, true, false)
}

var PostConsensusWrongValidatorIndexAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	msg := postConsensusAggregatorMsg(sk, id, true, false)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, false, true)
}

var postConsensusAggregatorMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainAggregateAndProof)

	aggData := TestingAggregateAndProof
	if wrongRoot {
		aggData = TestingWrongAggregateAndProof
	}

	signed, root, _ := signer.SignBeaconObject(aggData, d, sk.GetPublicKey().Serialize(), types.DomainAggregateAndProof)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(aggData, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainAggregateAndProof)
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed,
				SigningRoot:      root,
				Signer:           id,
				ValidatorIndex:   TestingValidatorIndex,
			},
		},
	}
	return &msgs
}

var PostConsensusSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID) *types.PartialSignatureMessages {
	return PostConsensusSyncCommitteeMsgForKeySetWithSlot(keySetMap, id, TestingDutySlot)
}

var PostConsensusSyncCommitteeMsgForKeySetWithSlot = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {

	var ret *types.PartialSignatureMessages
	// Get post consensus for sync committees for each validator in shares
	for valIdx, ks := range keySetMap {
		pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, slot, false, false, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}
	}
	return ret
}

var PostConsensusPartiallyWrongBeaconSigSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongSyncCommitteeMsgForKeySet(keySetMap, id, false, true)
}

var PostConsensusPartiallyWrongRootSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID) *types.PartialSignatureMessages {
	return PostConsensusPartiallyWrongSyncCommitteeMsgForKeySet(keySetMap, id, true, false)
}

var PostConsensusPartiallyWrongSyncCommitteeMsgForKeySet = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, id types.OperatorID, wrongRoot bool, wrongBeaconSig bool) *types.PartialSignatureMessages {

	numValid := len(keySetMap) / 2
	msgIndex := 0

	var ret *types.PartialSignatureMessages

	validatorIndexes := make([]phase0.ValidatorIndex, 0)
	for valIdx := range keySetMap {
		validatorIndexes = append(validatorIndexes, valIdx)
	}
	sort.Slice(validatorIndexes, func(i, j int) bool {
		return validatorIndexes[i] < validatorIndexes[j]
	})

	for _, valIdx := range validatorIndexes {
		ks, ok := keySetMap[valIdx]
		if !ok {
			panic("validator index not in key set map")
		}

		invalidMsgFlag := (msgIndex < numValid)

		wrongRootV := wrongRoot && invalidMsgFlag
		wrongBeaconSigV := wrongBeaconSig && invalidMsgFlag

		pSigMsgs := postConsensusSyncCommitteeMsg(ks.Shares[id], id, TestingDutySlot, wrongRootV, wrongBeaconSigV, valIdx)
		if ret == nil {
			ret = pSigMsgs
		} else {
			ret.Messages = append(ret.Messages, pSigMsgs.Messages...)
		}

		msgIndex++
	}
	return ret
}

var PostConsensusSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeMsg(sk, id, TestingDutySlot, false, false, TestingValidatorIndex)
}

var PostConsensusSyncCommitteeTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	ret := postConsensusSyncCommitteeMsg(sk, id, TestingDutySlot, false, false, TestingValidatorIndex)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusSyncCommitteeTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var PostConsensusWrongSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeMsg(sk, id, TestingDutySlot, true, false, TestingValidatorIndex)
}

var PostConsensusWrongValidatorIndexSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	msg := postConsensusSyncCommitteeMsg(sk, id, TestingDutySlot, true, false, TestingValidatorIndex)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeMsg(sk, id, TestingDutySlot, false, true, TestingValidatorIndex)
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
	d, _ := beacon.DomainData(1, types.DomainSyncCommittee)
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

var PreConsensusContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return PreConsensusCustomSlotContributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot)
}

var PreConsensusContributionProofWrongBeaconSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, false, true)
}

var PreConsensusContributionProofNextEpochMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot2, TestingDutySlot2, false, false)
}

var PreConsensusCustomSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, slot, TestingDutySlot, false, false)
}

var PreConsensusWrongMsgSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, false, false)
}

var PreConsensusWrongOrderContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, true, false)
}

var PreConsensusContributionProofTooManyRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	ret := contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: append(ret.Messages, ret.Messages[0]),
	}
	return msg
}

var PreConsensusContributionProofTooFewRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	ret := contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: ret.Messages[0:2],
	}
	return msg
}

var contributionProofMsg = func(
	sk, beaconsk *bls.SecretKey,
	id, beaconid types.OperatorID,
	slot phase0.Slot,
	msgSlot phase0.Slot,
	wrongMsgOrder bool,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainSyncCommitteeSelectionProof)

	msgs := make([]*types.PartialSignatureMessage, 0)
	for index := range TestingContributionProofIndexes {
		subnet, _ := beacon.SyncCommitteeSubnetID(phase0.CommitteeIndex(index))
		data := &altair.SyncAggregatorSelectionData{
			Slot:              slot,
			SubcommitteeIndex: subnet,
		}
		sig, root, _ := signer.SignBeaconObject(data, d, beaconsk.GetPublicKey().Serialize(), types.DomainSyncCommitteeSelectionProof)
		if wrongBeaconSig {
			sig, root, _ = signer.SignBeaconObject(data, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainSyncCommitteeSelectionProof)
		}

		msg := &types.PartialSignatureMessage{
			PartialSignature: sig[:],
			SigningRoot:      ensureRoot(root),
			Signer:           beaconid,
			ValidatorIndex:   TestingValidatorIndex,
		}

		msgs = append(msgs, msg)
	}

	if wrongMsgOrder {
		m := msgs[0]
		msgs[0] = msgs[1]
		msgs[1] = m
	}

	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: msgs,
	}
	return msg
}

var PostConsensusSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, false)
}

var PostConsensusSyncCommitteeContributionWrongOrderMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, true)
}

var PostConsensusSyncCommitteeContributionTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	ret := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, false)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusSyncCommitteeContributionTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	ret := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages[0:2],
	}

	return msg
}

var PostConsensusWrongSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, true, false, false)
}

var PostConsensusWrongValidatorIndexSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	msg := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, true, false, false)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, true, false)
}

var postConsensusSyncCommitteeContributionMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	validatorIndex phase0.ValidatorIndex,
	keySet *TestKeySet,
	wrongRoot bool,
	wrongBeaconSig bool,
	wrongRootOrder bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	dContribAndProof, _ := beacon.DomainData(1, types.DomainContributionAndProof)

	msgs := make([]*types.PartialSignatureMessage, 0)
	for index := range TestingSyncCommitteeContributions {
		// sign contrib and proof
		contribAndProof := &altair.ContributionAndProof{
			AggregatorIndex: validatorIndex,
			Contribution:    &TestingContributionsData[index].Contribution,
			SelectionProof:  TestingContributionsData[index].SelectionProofSig,
		}

		if wrongRoot {
			contribAndProof.AggregatorIndex = 100
		}

		signed, root, _ := signer.SignBeaconObject(contribAndProof, dContribAndProof, sk.GetPublicKey().Serialize(), types.DomainSyncCommitteeSelectionProof)
		if wrongBeaconSig {
			signed, root, _ = signer.SignBeaconObject(contribAndProof, dContribAndProof, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainSyncCommitteeSelectionProof)
		}

		msg := &types.PartialSignatureMessage{
			PartialSignature: signed,
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}

		msgs = append(msgs, msg)
	}

	if wrongRootOrder {
		m := msgs[0]
		msgs[0] = msgs[1]
		msgs[1] = m
	}

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: msgs,
	}

	return msg
}

// ensureRoot ensures that SigningRoot will have sufficient allocated memory
// otherwise we get panic from bls:
// github.com/herumi/bls-eth-go-binary/bls.(*Sign).VerifyByte:738
func ensureRoot(root [32]byte) [32]byte {
	tmp := [32]byte{}
	copy(tmp[:], root[:])
	return tmp
}
