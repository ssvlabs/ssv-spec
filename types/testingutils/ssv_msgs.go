package testingutils

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

var TestingSSVDomainType = types.JatoTestnet
var TestingForkData = types.ForkData{Epoch: TestingDutyEpoch, Domain: TestingSSVDomainType}
var AttesterMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleAttester)
	return ret[:]
}()

var ProposerMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleProposer)
	return ret[:]
}()
var AggregatorMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleAggregator)
	return ret[:]
}()
var SyncCommitteeMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleSyncCommittee)
	return ret[:]
}()
var SyncCommitteeContributionMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleSyncCommitteeContribution)
	return ret[:]
}()
var ValidatorRegistrationMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleValidatorRegistration)
	return ret[:]
}()
var VoluntaryExitMsgID = func() []byte {
	ret := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleVoluntaryExit)
	return ret[:]
}()

var TestAttesterConsensusData = &types.ConsensusData{
	Duty:    TestingAttesterDuty,
	DataSSZ: TestingAttestationDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

var TestAttesterNextEpochConsensusData = &types.ConsensusData{
	Duty:    TestingAttesterDutyNextEpoch,
	DataSSZ: TestingAttestationNextEpochDataBytes,
	Version: spec.DataVersionPhase0,
}

var TestingAttesterNextEpochConsensusDataByts, _ = TestAttesterNextEpochConsensusData.Encode()

var TestAggregatorConsensusData = &types.ConsensusData{
	Duty:    TestingAggregatorDuty,
	DataSSZ: TestingAggregateAndProofBytes,
	Version: spec.DataVersionPhase0,
}
var TestAggregatorConsensusDataByts, _ = TestAggregatorConsensusData.Encode()

var TestAttesterWithJustificationsConsensusData = func(ks *TestKeySet) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(1); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusRandaoMsg(ks.Shares[i], i))
	}

	return &types.ConsensusData{
		Duty:                       TestingAttesterDuty,
		Version:                    spec.DataVersionDeneb,
		PreConsensusJustifications: justif,
		DataSSZ:                    TestingAttestationDataBytes,
	}
}

var TestAggregatorWithJustificationsConsensusData = func(ks *TestKeySet) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(1); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusSelectionProofMsg(ks.Shares[i], ks.Shares[i], i, i))
	}

	return &types.ConsensusData{
		Duty:                       TestingAggregatorDuty,
		Version:                    spec.DataVersionBellatrix,
		PreConsensusJustifications: justif,
		DataSSZ:                    TestingAggregateAndProofBytes,
	}

}

// TestSyncCommitteeWithJustificationsConsensusData is an invalid sync committee msg (doesn't have pre-consensus)
var TestSyncCommitteeWithJustificationsConsensusData = func(ks *TestKeySet) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusRandaoMsg(ks.Shares[i+1], i+1))
	}

	return &types.ConsensusData{
		Duty:                       TestingSyncCommitteeDuty,
		Version:                    spec.DataVersionDeneb,
		PreConsensusJustifications: justif,
		DataSSZ:                    TestingSyncCommitteeBlockRoot[:],
	}
}

var TestSyncCommitteeContributionWithJustificationConsensusData = func(ks *TestKeySet) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusContributionProofMsg(ks.Shares[i+1], ks.Shares[i+1], i+1, i+1))
	}

	return &types.ConsensusData{
		Duty:                       TestingSyncCommitteeContributionDuty,
		Version:                    spec.DataVersionDeneb,
		PreConsensusJustifications: justif,
		DataSSZ:                    TestingContributionsDataBytes,
	}
}

var TestSyncCommitteeConsensusData = &types.ConsensusData{
	Duty:    TestingSyncCommitteeDuty,
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeConsensusDataByts, _ = TestSyncCommitteeConsensusData.Encode()

var TestSyncCommitteeNextEpochConsensusData = &types.ConsensusData{
	Duty:    TestingSyncCommitteeDutyNextEpoch,
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: spec.DataVersionPhase0,
}

var TestSyncCommitteeNextEpochConsensusDataByts, _ = TestSyncCommitteeNextEpochConsensusData.Encode()

var TestSyncCommitteeContributionConsensusData = &types.ConsensusData{
	Duty:    TestingSyncCommitteeContributionDuty,
	DataSSZ: TestingContributionsDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()
var TestSyncCommitteeContributionConsensusDataRoot = func() [32]byte {
	return sha256.Sum256(TestSyncCommitteeContributionConsensusDataByts)
}()

var TestConsensusUnkownDutyTypeData = &types.ConsensusData{
	Duty:    TestingUnknownDutyType,
	DataSSZ: TestingAttestationDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var TestConsensusWrongDutyPKData = &types.ConsensusData{
	Duty:    TestingWrongDutyPK,
	DataSSZ: TestingAttestationDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestConsensusWrongDutyPKDataByts, _ = TestConsensusWrongDutyPKData.Encode()

var SSVMsgAttester = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleAttester))
}

var SSVMsgWrongID = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingWrongValidatorPubKey[:], types.BNRoleAttester))
}

var SSVMsgProposer = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleProposer))
}

var SSVMsgAggregator = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleAggregator))
}

var SSVMsgSyncCommittee = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleSyncCommittee))
}

var SSVMsgSyncCommitteeContribution = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleSyncCommitteeContribution))
}

var SSVMsgValidatorRegistration = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleValidatorRegistration))
}

var SSVMsgVoluntaryExit = func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleVoluntaryExit))
}

var ssvMsg = func(qbftMsg *qbft.SignedMessage, postMsg *types.SignedPartialSignatureMessage, msgID types.MessageID) *types.SSVMessage {
	var msgType types.MsgType
	var data []byte
	var err error
	if qbftMsg != nil {
		msgType = types.SSVConsensusMsgType
		data, err = qbftMsg.Encode()
		if err != nil {
			panic(err)
		}
	} else if postMsg != nil {
		msgType = types.SSVPartialSignatureMsgType
		data, err = postMsg.Encode()
		if err != nil {
			panic(err)
		}
	} else {
		panic("msg type undefined")
	}

	return &types.SSVMessage{
		MsgType: msgType,
		MsgID:   msgID,
		Data:    data,
	}
}

var PostConsensusWrongAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.SignedPartialSignatureMessage {
	return postConsensusAttestationMsg(sk, id, height, true, false)
}

var PostConsensusWrongSigAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.SignedPartialSignatureMessage {
	return postConsensusAttestationMsg(sk, id, height, false, true)
}

var PostConsensusSigAttestationWrongBeaconSignerMsg = func(sk *bls.SecretKey, id, beaconSigner types.OperatorID, height qbft.Height) *types.SignedPartialSignatureMessage {
	ret := postConsensusAttestationMsg(sk, beaconSigner, height, false, true)
	ret.Signer = id
	return ret
}

var PostConsensusAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.SignedPartialSignatureMessage {
	return postConsensusAttestationMsg(sk, id, height, false, false)
}

var PostConsensusAttestationTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.SignedPartialSignatureMessage {
	ret := postConsensusAttestationMsg(sk, id, height, false, false)
	ret.Message.Messages = append(ret.Message.Messages, ret.Message.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages,
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusAttestationTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *types.SignedPartialSignatureMessage {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var postConsensusAttestationMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	height qbft.Height,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingAttestationData.Target.Epoch, types.DomainAttester)

	attData := TestingAttestationData
	if wrongRoot {
		attData = TestingWrongAttestationData
	}

	signed, root, _ := signer.SignBeaconObject(attData, d, sk.GetPublicKey().Serialize(), types.DomainAttester)

	if wrongBeaconSig {
		signed, _, _ = signer.SignBeaconObject(attData, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainAttester)
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed,
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusFailedMsg = func(msgSigner *bls.SecretKey, msgSignerID types.OperatorID) *types.SignedPartialSignatureMessage {
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
			},
		},
	}
	sig, _ := signer.SignRoot(msg, types.PartialSignatureType, msgSigner.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msg,
		Signature: sig,
		Signer:    msgSignerID,
	}
}

var PreConsensusRandaoMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 1, false)
}

var randaoMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	epoch phase0.Epoch,
	msgCnt int,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
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
		}
		if wrongRoot {
			msg.SigningRoot = [32]byte{}
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return PreConsensusCustomSlotSelectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot)
}

var PreConsensusSelectionProofWrongBeaconSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, 1, true)
}

var PreConsensusSelectionProofNextEpochMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot2, TestingDutySlot2, 1, false)
}

var PreConsensusSelectionProofTooManyRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, 3, false)
}

var PreConsensusSelectionProofTooFewRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, 0, false)
}

var PreConsensusCustomSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, slot, TestingDutySlot, 1, false)
}

var PreConsensusWrongMsgSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, 1, false)
}

var TestSelectionProofWithJustificationsConsensusData = func(ks *TestKeySet) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusSelectionProofMsg(ks.Shares[i+1], ks.Shares[i+1], i+1, i+1))
	}

	return &types.ConsensusData{
		Duty:                       TestingAggregatorDuty,
		Version:                    spec.DataVersionDeneb,
		PreConsensusJustifications: justif,
		DataSSZ:                    TestingAggregateAndProofBytes,
	}
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
) *types.SignedPartialSignatureMessage {
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
		})
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.SelectionProofPartialSig,
		Slot:     TestingDutySlot,
		Messages: _msgs,
	}
	msgSig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: msgSig,
		Signer:    id,
	}
}

var PreConsensusValidatorRegistrationMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var validatorRegistrationMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
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
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	msg := &types.PartialSignatureMessage{
		PartialSignature: signed[:],
		SigningRoot:      root,
		Signer:           id,
	}
	if wrongRoot {
		msg.SigningRoot = [32]byte{}
	}

	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusVoluntaryExitMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var PreConsensusVoluntaryExitTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var VoluntaryExitMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
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
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	msg := &types.PartialSignatureMessage{
		PartialSignature: signed[:],
		SigningRoot:      root,
		Signer:           id,
	}
	if wrongRoot {
		msg.SigningRoot = [32]byte{}
	}

	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusAggregatorMsg(sk, id, false, false)
}

var PostConsensusAggregatorTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := postConsensusAggregatorMsg(sk, id, false, false)
	ret.Message.Messages = append(ret.Message.Messages, ret.Message.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages,
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusAggregatorTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusWrongAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusAggregatorMsg(sk, id, true, false)
}

var PostConsensusWrongSigAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusAggregatorMsg(sk, id, false, true)
}

var PostConsensusSigAggregatorWrongBeaconSignerMsg = func(sk *bls.SecretKey, id, beaconSigner types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := postConsensusAggregatorMsg(sk, beaconSigner, false, true)
	ret.Signer = id
	return ret
}

var postConsensusAggregatorMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
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
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeMsg(sk, id, false, false)
}

var PostConsensusSyncCommitteeTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := postConsensusSyncCommitteeMsg(sk, id, false, false)
	ret.Message.Messages = append(ret.Message.Messages, ret.Message.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages,
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusSyncCommitteeTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusWrongSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeMsg(sk, id, true, false)
}

var PostConsensusWrongSigSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeMsg(sk, id, false, true)
}

var PostConsensusSigSyncCommitteeWrongBeaconSignerMsg = func(sk *bls.SecretKey, id, beaconSigner types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := postConsensusSyncCommitteeMsg(sk, beaconSigner, false, true)
	ret.Signer = id
	return ret
}

var postConsensusSyncCommitteeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
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
		Slot: TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed,
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return PreConsensusCustomSlotContributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot)
}

var PreConsensusContributionProofWrongBeaconSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, false, true)
}

var PreConsensusContributionProofNextEpochMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot2, TestingDutySlot2, false, false)
}

var PreConsensusCustomSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, slot, TestingDutySlot, false, false)
}

var PreConsensusWrongMsgSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, false, false)
}

var PreConsensusWrongOrderContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, true, false)
}
var TestContributionProofWithJustificationsConsensusData = func(ks *TestKeySet) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusContributionProofMsg(ks.Shares[i+1], ks.Shares[i+1], i+1, i+1))
	}

	return &types.ConsensusData{
		Duty:                       TestingSyncCommitteeContributionDuty,
		Version:                    spec.DataVersionDeneb,
		PreConsensusJustifications: justif,
		DataSSZ:                    TestingContributionsDataBytes,
	}
}

var PreConsensusContributionProofTooManyRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: append(ret.Message.Messages, ret.Message.Messages[0]),
	}

	msgSig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, beaconSK.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: msgSig,
		Signer:    msgID,
	}
}

var PreConsensusContributionProofTooFewRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages[0:2],
	}

	msgSig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, beaconSK.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: msgSig,
		Signer:    msgID,
	}
}

var contributionProofMsg = func(
	sk, beaconsk *bls.SecretKey,
	id, beaconid types.OperatorID,
	slot phase0.Slot,
	msgSlot phase0.Slot,
	wrongMsgOrder bool,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
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

	msgSig, _ := signer.SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: msgSig,
		Signer:    id,
	}
}

var PostConsensusSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, false)
}

var PostConsensusSyncCommitteeContributionWrongOrderMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, true)
}

var PostConsensusSyncCommitteeContributionTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	ret := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, false)
	ret.Message.Messages = append(ret.Message.Messages, ret.Message.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages,
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusSyncCommitteeContributionTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	ret := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages[0:2],
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusWrongSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, true, false, false)
}

var PostConsensusWrongSigSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, true, false)
}

var PostConsensusSigSyncCommitteeContributionWrongSignerMsg = func(sk *bls.SecretKey, id, beaconSigner types.OperatorID, keySet *TestKeySet) *types.SignedPartialSignatureMessage {
	ret := postConsensusSyncCommitteeContributionMsg(sk, beaconSigner, TestingValidatorIndex, keySet, false, false, false)
	ret.Signer = id
	return ret
}

var postConsensusSyncCommitteeContributionMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	validatorIndex phase0.ValidatorIndex,
	keySet *TestKeySet,
	wrongRoot bool,
	wrongBeaconSig bool,
	wrongRootOrder bool,
) *types.SignedPartialSignatureMessage {
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

	sig, _ := signer.SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

// ensureRoot ensures that SigningRoot will have sufficient allocated memory
// otherwise we get panic from bls:
// github.com/herumi/bls-eth-go-binary/bls.(*Sign).VerifyByte:738
func ensureRoot(root [32]byte) [32]byte {
	tmp := [32]byte{}
	copy(tmp[:], root[:])
	return tmp
}
