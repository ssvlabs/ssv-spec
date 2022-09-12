package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var AttesterMsgID = func() []byte {
	ret := types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleAttester)
	return ret[:]
}()

var ProposerMsgID = func() []byte {
	ret := types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleProposer)
	return ret[:]
}()
var AggregatorMsgID = func() []byte {
	ret := types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleAggregator)
	return ret[:]
}()
var SyncCommitteeMsgID = func() []byte {
	ret := types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleSyncCommittee)
	return ret[:]
}()
var SyncCommitteeContributionMsgID = func() []byte {
	ret := types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleSyncCommitteeContribution)
	return ret[:]
}()

var TestAttesterConsensusData = &types.ConsensusData{
	Duty:            TestingAttesterDuty,
	AttestationData: TestingAttestationData,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

var TestAggregatorConsensusData = &types.ConsensusData{
	Duty:              TestingAggregatorDuty,
	AggregateAndProof: TestingAggregateAndProof,
}
var TestAggregatorConsensusDataByts, _ = TestAggregatorConsensusData.Encode()

var TestProposerConsensusData = &types.ConsensusData{
	Duty:      TestingProposerDuty,
	BlockData: TestingBeaconBlock,
}
var TestProposerConsensusDataByts, _ = TestProposerConsensusData.Encode()

var TestSyncCommitteeConsensusData = &types.ConsensusData{
	Duty:                   TestingSyncCommitteeDuty,
	SyncCommitteeBlockRoot: TestingSyncCommitteeBlockRoot,
}
var TestSyncCommitteeConsensusDataByts, _ = TestSyncCommitteeConsensusData.Encode()

var TestSyncCommitteeContributionConsensusData = &types.ConsensusData{
	Duty: TestingSyncCommitteeContributionDuty,
	SyncCommitteeContribution: map[spec.BLSSignature]*altair.SyncCommitteeContribution{
		TestingContributionProofsSigned[0]: TestingSyncCommitteeContributions[0],
		TestingContributionProofsSigned[1]: TestingSyncCommitteeContributions[1],
		TestingContributionProofsSigned[2]: TestingSyncCommitteeContributions[2],
	},
}
var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()

var TestConsensusUnkownDutyTypeData = &types.ConsensusData{
	Duty:            TestingUnknownDutyType,
	AttestationData: TestingAttestationData,
}
var TestConsensusUnkownDutyTypeDataByts, _ = TestConsensusUnkownDutyTypeData.Encode()

var TestConsensusWrongDutyPKData = &types.ConsensusData{
	Duty:            TestingWrongDutyPK,
	AttestationData: TestingAttestationData,
}
var TestConsensusWrongDutyPKDataByts, _ = TestConsensusWrongDutyPKData.Encode()

var SSVMsgAttester = func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleAttester))
}

var SSVMsgWrongID = func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingWrongValidatorPubKey[:], types.BNRoleAttester))
}

var SSVMsgProposer = func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleProposer))
}

var SSVMsgAggregator = func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleAggregator))
}

var SSVMsgSyncCommittee = func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleSyncCommittee))
}

var SSVMsgSyncCommitteeContribution = func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingValidatorPubKey[:], types.BNRoleSyncCommitteeContribution))
}

var ssvMsg = func(qbftMsg *qbft.SignedMessage, postMsg *ssv.SignedPartialSignatureMessage, msgID types.MessageIDOld) *types.SSVMessage {
	var msgType types.MsgType
	var data []byte
	if qbftMsg != nil {
		//msgType = types.SSVConsensusMsgType
		data, _ = qbftMsg.Encode()
	} else if postMsg != nil {
		//msgType = types.SSVPartialSignatureMsgType
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

var PostConsensusAttestationMsgWithWrongSig = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPartialSignatureMessage {
	return postConsensusAttestationMsg(sk, id, height, true, false)
}

var PostConsensusAttestationMsgWithWrongRoot = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPartialSignatureMessage {
	return postConsensusAttestationMsg(sk, id, height, true, false)
}

var PostConsensusAttestationMsg = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *ssv.SignedPartialSignatureMessage {
	return postConsensusAttestationMsg(sk, id, height, false, false)
}

var postConsensusAttestationMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	height qbft.Height,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	signed, root, _ := signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, sk.GetPublicKey().Serialize())

	if wrongBeaconSig {
		signed, _, _ = signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, TestingWrongValidatorPubKey[:])
	}

	if wrongRoot {
		root = []byte{1, 2, 3, 4}
	}

	msgs := ssv.PartialSignatureMessages{
		Type: ssv.PostConsensusPartialSig,
		Messages: []*ssv.PartialSignatureMessage{
			{
				Slot:             TestingDutySlot,
				PartialSignature: signed.Signature[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusProposerMsg = func(sk *bls.SecretKey, id types.OperatorID) *ssv.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsg(sk, id, false, false)
}

var postConsensusBeaconBlockMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	signed, root, _ := signer.SignBeaconBlock(TestingBeaconBlock, TestingProposerDuty, sk.GetPublicKey().Serialize())

	if wrongBeaconSig {
		//signed, _, _ = signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, TestingWrongSK.GetPublicKey().Serialize())
		panic("implement")
	}

	if wrongRoot {
		root = []byte{1, 2, 3, 4}
	}

	msgs := ssv.PartialSignatureMessages{
		Type: ssv.PostConsensusPartialSig,
		Messages: []*ssv.PartialSignatureMessage{
			{
				Slot:             TestingDutySlot,
				PartialSignature: signed.Signature[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusRandaoMsg = func(sk *bls.SecretKey, id types.OperatorID) *ssv.SignedPartialSignatureMessage {
	return randaoMsg(sk, id, false, false)
}

var randaoMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	signed, root, _ := signer.SignRandaoReveal(TestingDutySlot, sk.GetPublicKey().Serialize())

	msgs := ssv.PartialSignatureMessages{
		Type: ssv.RandaoPartialSig,
		Messages: []*ssv.PartialSignatureMessage{
			{
				Slot:             TestingDutySlot,
				PartialSignature: signed[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusSelectionProofMsg = func(sk *bls.SecretKey, id types.OperatorID) *ssv.SignedPartialSignatureMessage {
	return selectionProofMsg(sk, id, false, false)
}

var selectionProofMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	signed, root, _ := signer.SignSlotWithSelectionProof(TestingDutySlot, sk.GetPublicKey().Serialize())

	msgs := ssv.PartialSignatureMessages{
		Type: ssv.SelectionProofPartialSig,
		Messages: []*ssv.PartialSignatureMessage{
			{
				Slot:             TestingDutySlot,
				PartialSignature: signed[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	msgSig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: msgSig,
		Signer:    id,
	}
}

var PostConsensusAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID) *ssv.SignedPartialSignatureMessage {
	return postConsensusAggregatorMsg(sk, id, false, false)
}

var postConsensusAggregatorMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	signed, root, _ := signer.SignAggregateAndProof(TestingAggregateAndProof, TestingProposerDuty, sk.GetPublicKey().Serialize())

	if wrongBeaconSig {
		//signed, _, _ = signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, TestingWrongSK.GetPublicKey().Serialize())
		panic("implement")
	}

	if wrongRoot {
		root = []byte{1, 2, 3, 4}
	}

	msgs := ssv.PartialSignatureMessages{
		Type: ssv.PostConsensusPartialSig,
		Messages: []*ssv.PartialSignatureMessage{
			{
				Slot:             TestingDutySlot,
				PartialSignature: signed.Signature[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusSyncCommitteeMsg = func(sk *bls.SecretKey, id types.OperatorID) *ssv.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeMsg(sk, id, false, false)
}

var postConsensusSyncCommitteeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	signed, root, _ := signer.SignSyncCommitteeBlockRoot(TestingDutySlot, TestingSyncCommitteeBlockRoot, TestingSyncCommitteeDuty.ValidatorIndex, sk.GetPublicKey().Serialize())

	if wrongBeaconSig {
		//signedAtt, _, _ = signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, TestingWrongSK.GetPublicKey().Serialize())
		panic("implement")
	}

	if wrongRoot {
		root = []byte{1, 2, 3, 4}
	}

	msgs := ssv.PartialSignatureMessages{
		Type: ssv.PostConsensusPartialSig,
		Messages: []*ssv.PartialSignatureMessage{
			{
				Slot:             TestingDutySlot,
				PartialSignature: signed.Signature[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}

var PreConsensusContributionProofMsg = func(sk *bls.SecretKey, id types.OperatorID) *ssv.SignedPartialSignatureMessage {
	return contributionProofMsg(sk, id, false, false)
}

var contributionProofMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	msgs := make([]*ssv.PartialSignatureMessage, 0)
	for index, _ := range TestingContributionProofRoots {
		sig, root, _ := signer.SignContributionProof(TestingDutySlot, uint64(index), sk.GetPublicKey().Serialize())
		msg := &ssv.PartialSignatureMessage{
			Slot:             TestingDutySlot,
			PartialSignature: sig[:],
			SigningRoot:      root,
			Signer:           id,
			MetaData: &ssv.PartialSignatureMetaData{
				ContributionSubCommitteeIndex: uint64(index),
			},
		}
		msgs = append(msgs, msg)
	}

	msg := &ssv.PartialSignatureMessages{
		Type:     ssv.ContributionProofs,
		Messages: msgs,
	}

	msgSig, _ := signer.SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: msgSig,
		Signer:    id,
	}
}

var PostConsensusSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *ssv.SignedPartialSignatureMessage {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, keySet, false, false)
}

var postConsensusSyncCommitteeContributionMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	validatorIndex spec.ValidatorIndex,
	keySet *TestKeySet,
	wrongRoot bool,
	wrongBeaconSig bool,
) *ssv.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()

	msgs := make([]*ssv.PartialSignatureMessage, 0)
	for index, c := range TestingSyncCommitteeContributions {
		signedProof, _, _ := signer.SignContributionProof(TestingDutySlot, uint64(index), keySet.ValidatorSK.GetPublicKey().Serialize())
		signedProofbls := spec.BLSSignature{}
		copy(signedProofbls[:], signedProof)

		signed, root, _ := signer.SignContribution(&altair.ContributionAndProof{
			AggregatorIndex: validatorIndex,
			Contribution:    c,
			SelectionProof:  signedProofbls,
		}, sk.GetPublicKey().Serialize())

		if wrongRoot {
			root = []byte{1, 2, 3, 4}
		}

		msg := &ssv.PartialSignatureMessage{
			Slot:             TestingDutySlot,
			PartialSignature: signed.Signature[:],
			SigningRoot:      root,
			Signer:           id,
		}

		if wrongBeaconSig {
			//signedAtt, _, _ = signer.SignAttestation(TestingAttestationData, TestingAttesterDuty, TestingWrongSK.GetPublicKey().Serialize())
			panic("implement")
		}

		msgs = append(msgs, msg)
	}

	msg := &ssv.PartialSignatureMessages{
		Type:     ssv.PostConsensusPartialSig,
		Messages: msgs,
	}

	sig, _ := signer.SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &ssv.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}
