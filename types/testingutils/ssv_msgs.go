package testingutils

import (
	spec2 "github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingSSVDomainType = types.V3Testnet
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

var TestAttesterConsensusData = &types.ConsensusData{
	Duty:    TestingAttesterDuty,
	DataSSZ: TestingAttestationDataBytes,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

var TestAggregatorConsensusData = &types.ConsensusData{
	Duty:    TestingAggregatorDuty,
	DataSSZ: TestingAggregateAndProofBytes,
}
var TestAggregatorConsensusDataByts, _ = TestAggregatorConsensusData.Encode()

var TestProposerConsensusData = &types.ConsensusData{
	Duty:    TestingProposerDuty,
	Version: spec2.DataVersionBellatrix,
	DataSSZ: TestingBeaconBlockBytes,
}
var TestProposerConsensusDataByts, _ = TestProposerConsensusData.Encode()

var TestProposerBlindedBlockConsensusData = &types.ConsensusData{
	Duty:    TestingProposerDuty,
	Version: spec2.DataVersionBellatrix,
	DataSSZ: TestingBlindedBeaconBlockBytes,
}
var TestProposerBlindedBlockConsensusDataByts, _ = TestProposerBlindedBlockConsensusData.Encode()

var TestSyncCommitteeConsensusData = &types.ConsensusData{
	Duty:    TestingSyncCommitteeDuty,
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
}
var TestSyncCommitteeConsensusDataByts, _ = TestSyncCommitteeConsensusData.Encode()

var TestSyncCommitteeContributionConsensusData = &types.ConsensusData{
	Duty:    TestingSyncCommitteeContributionDuty,
	DataSSZ: TestingSyncCommitteeContributionsConsensusDataBytes,
}
var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()

var TestConsensusUnkownDutyTypeData = &types.ConsensusData{
	Duty:    TestingUnknownDutyType,
	DataSSZ: TestingAttestationDataBytes,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var TestConsensusWrongDutyPKData = &types.ConsensusData{
	Duty:    TestingWrongDutyPK,
	DataSSZ: TestingAttestationDataBytes,
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

var ssvMsg = func(qbftMsg *qbft.SignedMessage, postMsg *types.SignedPartialSignatureMessage, msgID types.MessageID) *types.SSVMessage {
	var msgType types.MsgType
	var data []byte
	if qbftMsg != nil {
		msgType = types.SSVConsensusMsgType
		data, _ = qbftMsg.Encode()
	} else if postMsg != nil {
		msgType = types.SSVPartialSignatureMsgType
		data, _ = postMsg.Encode()
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

var PostConsensusProposerMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsg(sk, id, false, false)
}

var PostConsensusProposerTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := postConsensusBeaconBlockMsg(sk, id, false, false)
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

var PostConsensusProposerTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
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

var PostConsensusWrongProposerMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsg(sk, id, true, false)
}

var PostConsensusWrongSigProposerMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsg(sk, id, false, true)
}

var PostConsensusSigProposerWrongBeaconSignerMsg = func(sk *bls.SecretKey, id, beaconSigner types.OperatorID) *types.SignedPartialSignatureMessage {
	ret := postConsensusBeaconBlockMsg(sk, beaconSigner, false, true)
	ret.Signer = id
	return ret
}

var postConsensusBeaconBlockMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()

	block := TestingBeaconBlock
	if wrongRoot {
		block = TestingWrongBeaconBlock
	}

	d, _ := beacon.DomainData(1, types.DomainProposer) // epoch doesn't matter here, hard coded
	sig, root, _ := signer.SignBeaconObject(block, d, sk.GetPublicKey().Serialize(), types.DomainProposer)
	if wrongBeaconSig {
		sig, root, _ = signer.SignBeaconObject(block, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainProposer)
	}
	blsSig := spec.BLSSignature{}
	copy(blsSig[:], sig)

	signed := bellatrix.SignedBeaconBlock{
		Message:   TestingBeaconBlock,
		Signature: blsSig,
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed.Signature[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	msgSig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: msgSig,
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

// PreConsensusRandaoNextEpochMsg testing for a second duty start
var PreConsensusRandaoNextEpochMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch2, 1, false)
}

var PreConsensusRandaoDifferentEpochMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch+1, 1, false)
}

var PreConsensusRandaoTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 2, false)
}

var PreConsensusRandaoTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 0, false)
}

var PreConsensusRandaoNoMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 0, false)
}

var PreConsensusRandaoWrongBeaconSigMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 1, true)
}

var PreConsensusRandaoDifferentSignerMsg = func(
	msgSigner, randaoSigner *bls.SecretKey,
	msgSignerID,
	randaoSignerID types.OperatorID,
) *types.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainRandao)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(TestingDutyEpoch), d, randaoSigner.GetPublicKey().Serialize(), types.DomainRandao)

	msg := types.PartialSignatureMessages{
		Type: types.RandaoPartialSig,
		Slot: TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed[:],
				SigningRoot:      root,
				Signer:           randaoSignerID,
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

var randaoMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	epoch spec.Epoch,
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

var PreConsensusCustomSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot spec.Slot) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, slot, TestingDutySlot, 1, false)
}

var PreConsensusWrongMsgSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return selectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, 1, false)
}

var selectionProofMsg = func(
	sk *bls.SecretKey,
	beaconsk *bls.SecretKey,
	id types.OperatorID,
	beaconid types.OperatorID,
	slot spec.Slot,
	msgSlot spec.Slot,
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
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutyEpoch, false)
}

var PreConsensusValidatorRegistrationTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutyEpoch, false)
}

var PreConsensusValidatorRegistrationTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutyEpoch, false)
}

var PreConsensusValidatorRegistrationDifferentEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.SignedPartialSignatureMessage {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutyEpoch, false)
}

var validatorRegistrationMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	epoch spec.Epoch,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(epoch, types.DomainApplicationBuilder)

	signed, root, _ := signer.SignBeaconObject(TestingValidatorRegistration, d, beaconSK.GetPublicKey().Serialize(), types.DomainApplicationBuilder)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingValidatorRegistrationWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainApplicationBuilder)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingValidatorRegistration, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainApplicationBuilder)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.ValidatorRegistrationPartialSig,
		Slot:     TestingDutySlot,
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

var PreConsensusCustomSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot spec.Slot) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, slot, TestingDutySlot, false, false)
}

var PreConsensusWrongMsgSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot+1, false, false)
}

var PreConsensusWrongOrderContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.SignedPartialSignatureMessage {
	return contributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot, TestingDutySlot, true, false)
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
	slot spec.Slot,
	msgSlot spec.Slot,
	wrongMsgOrder bool,
	wrongBeaconSig bool,
) *types.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainSyncCommitteeSelectionProof)

	msgs := make([]*types.PartialSignatureMessage, 0)
	for index := range TestingContributionProofIndexes {
		subnet, _ := beacon.SyncCommitteeSubnetID(spec.CommitteeIndex(index))
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
	ret := postConsensusSyncCommitteeContributionMsg(sk, beaconSigner, TestingValidatorIndex, keySet, false, true, false)
	ret.Signer = id
	return ret
}

var postConsensusSyncCommitteeContributionMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	validatorIndex spec.ValidatorIndex,
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
			Contribution:    &TestingSyncCommitteeContributionsConsensusData[index].Contribution,
			SelectionProof:  TestingSyncCommitteeContributionsConsensusData[index].SelectionProofSig,
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
