package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgSyncCommitteeContribution = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleSyncCommitteeContribution))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return PreConsensusCustomSlotContributionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlot)
}

var PreConsensusContributionProofWrongBeaconSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlot, false, true)
}

var PreConsensusContributionProofNextEpochMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlot2, false, false)
}

var PreConsensusCustomSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, slot, false, false)
}

var PreConsensusContributionProofMsgWithValidatorIndex = func(msgSK, beaconSK *bls.SecretKey, validatorIndex phase0.ValidatorIndex, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return PreConsensusCustomSlotContributionProofMsgWithValidatorIndex(msgSK, beaconSK, validatorIndex, msgID, beaconID, TestingDutySlot)
}

var PreConsensusCustomSlotContributionProofMsgWithValidatorIndex = func(msgSK, beaconSK *bls.SecretKey, validatorIndex phase0.ValidatorIndex, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, validatorIndex, msgID, beaconID, slot, false, false)
}

var PreConsensusWrongMsgSlotContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlot, false, false)
}

var PreConsensusWrongOrderContributionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	return contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlot, true, false)
}

var PreConsensusContributionProofTooManyRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	ret := contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlot, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: append(ret.Messages, ret.Messages[0]),
	}
	return msg
}

var PreConsensusContributionProofTooFewRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID) *types.PartialSignatureMessages {
	ret := contributionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlot, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     TestingDutySlot,
		Messages: ret.Messages[0:2],
	}
	return msg
}

var contributionProofMsg = func(
	sk, beaconsk *bls.SecretKey,
	validatorIndex phase0.ValidatorIndex,
	id, beaconid types.OperatorID,
	slot phase0.Slot,
	wrongMsgOrder bool,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainSyncCommitteeSelectionProof)

	msgs := make([]*types.PartialSignatureMessage, 0)
	for index := range TestingContributionProofIndexes {
		subnet := beacon.SyncCommitteeSubnetID(phase0.CommitteeIndex(index))
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
			ValidatorIndex:   validatorIndex,
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

// ==================================================
// PostConsensus
// ==================================================

var PostConsensusSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, false, false, false)
}

var PostConsensusSyncCommitteeContributionWrongOrderMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, false, false, true)
}

var PostConsensusSyncCommitteeContributionTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	ret := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, false, false, false)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusSyncCommitteeContributionTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	ret := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, false, false, false)
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlot,
		Messages: ret.Messages[0:2],
	}

	return msg
}

var PostConsensusWrongSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, true, false, false)
}

var PostConsensusWrongValidatorIndexSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	msg := postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, true, false, false)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigSyncCommitteeContributionMsg = func(sk *bls.SecretKey, id types.OperatorID, keySet *TestKeySet) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, TestingValidatorIndex, TestingDutySlot, keySet, false, true, false)
}

var PostConsensusCustomSlotSyncCommitteeContributionMsgWithValidatorIndex = func(sk *bls.SecretKey, validatorIndex phase0.ValidatorIndex, id types.OperatorID, keySet *TestKeySet, slot phase0.Slot) *types.PartialSignatureMessages {
	return postConsensusSyncCommitteeContributionMsg(sk, id, validatorIndex, slot, keySet, false, false, false)
}

var postConsensusSyncCommitteeContributionMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	validatorIndex phase0.ValidatorIndex,
	slot phase0.Slot,
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
			ValidatorIndex:   validatorIndex,
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
		Slot:     slot,
		Messages: msgs,
	}

	return msg
}
