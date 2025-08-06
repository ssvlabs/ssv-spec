package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgAggregator = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleAggregator))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PreConsensusCustomSlotSelectionProofMsg(msgSK, beaconSK, msgID, beaconID, TestingDutySlotV(version))
}

var PreConsensusSelectionProofWrongBeaconSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlotV(version), 1, true, false)
}

var PreConsensusSelectionProofWrongRootSigMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlotV(version), 1, false, true)
}

var PreConsensusSelectionProofNextEpochMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlotNextEpochV(version), 1, false, false)
}

var PreConsensusSelectionProofTooManyRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlotV(version), 3, false, false)
}

var PreConsensusSelectionProofTooFewRootsMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlotV(version), 0, false, false)
}

var PreConsensusCustomSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, slot, 1, false, false)
}

var PreConsensusSelectionProofMsgWithValidatorIndex = func(msgSK, beaconSK *bls.SecretKey, validatorIndex phase0.ValidatorIndex, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return PreConsensusCustomSlotSelectionProofMsgWithValidatorIndex(msgSK, beaconSK, validatorIndex, msgID, beaconID, TestingDutySlotV(version))
}

var PreConsensusCustomSlotSelectionProofMsgWithValidatorIndex = func(msgSK, beaconSK *bls.SecretKey, validatorIndex phase0.ValidatorIndex, msgID, beaconID types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, validatorIndex, msgID, beaconID, slot, 1, false, false)
}

var PreConsensusWrongMsgSlotSelectionProofMsg = func(msgSK, beaconSK *bls.SecretKey, msgID, beaconID types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return selectionProofMsg(msgSK, beaconSK, TestingValidatorIndex, msgID, beaconID, TestingDutySlotV(version), 1, false, false)
}

var selectionProofMsg = func(
	sk *bls.SecretKey,
	beaconsk *bls.SecretKey,
	validatorIndex phase0.ValidatorIndex,
	id types.OperatorID,
	beaconid types.OperatorID,
	slot phase0.Slot,
	msgCnt int,
	wrongBeaconSig bool,
	wrongSigningRoot bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainSelectionProof)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(slot), d, beaconsk.GetPublicKey().Serialize(), types.DomainSelectionProof)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(types.SSZUint64(slot), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainSelectionProof)
	}
	if wrongSigningRoot {
		_, root, _ = signer.SignBeaconObject(types.SSZUint64(slot+1), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainSelectionProof)
	}

	_msgs := make([]*types.PartialSignatureMessage, 0)
	for i := 0; i < msgCnt; i++ {
		_msgs = append(_msgs, &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           beaconid,
			ValidatorIndex:   validatorIndex,
		})
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.SelectionProofPartialSig,
		Slot:     slot,
		Messages: _msgs,
	}
	return &msgs
}

// ==================================================
// PostConsensus
// ==================================================

var PostConsensusAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, TestingValidatorIndex, false, false, version)
}

var PostConsensusAggregatorTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := postConsensusAggregatorMsg(sk, id, TestingValidatorIndex, false, false, version)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusAggregatorTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotV(version),
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var PostConsensusWrongAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, TestingValidatorIndex, true, false, version)
}

var PostConsensusWrongValidatorIndexAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusAggregatorMsg(sk, id, TestingValidatorIndex, true, false, version)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigAggregatorMsg = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, TestingValidatorIndex, false, true, version)
}

var PostConsensusAggregatorMsgWithValidatorIndex = func(sk *bls.SecretKey, id types.OperatorID, validatorIndex phase0.ValidatorIndex, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusAggregatorMsg(sk, id, validatorIndex, false, false, version)
}

var postConsensusAggregatorMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	validatorIndex phase0.ValidatorIndex,
	wrongRoot bool,
	wrongBeaconSig bool,
	version spec.DataVersion,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(1, types.DomainAggregateAndProof)

	aggData := TestingAggregateAndProofV(version, validatorIndex)
	if wrongRoot {
		aggData = TestingWrongAggregateAndProofV(version, validatorIndex)
	}

	signed, root, _ := signer.SignBeaconObject(aggData.(ssz.HashRoot), d, sk.GetPublicKey().Serialize(), types.DomainAggregateAndProof)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(aggData.(ssz.HashRoot), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainAggregateAndProof)
	}

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingDutySlotV(version),
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
