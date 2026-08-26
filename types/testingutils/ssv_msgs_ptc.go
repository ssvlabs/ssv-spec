package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgPTCAttester = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RolePTCAttester))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusPTCMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 1, false, TestingDutySlotGloas, false)
}

var PreConsensusPTCTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 0, false, TestingDutySlotGloas, false)
}

var PreConsensusPTCTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 2, false, TestingDutySlotGloas, false)
}

var PreConsensusPTCWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 1, false, TestingDutySlotGloas, true)
}

var PreConsensusPTCNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 1, false, TestingDutySlotGloasNextEpoch, false)
}

// PreConsensusPTCAbstainSlotMsg is a peer's valid partial signature for the abstain slot — the peer
// observed a block there even though the local beacon node reports none.
var PreConsensusPTCAbstainSlotMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 1, false, TestingPTCAbstainSlot, false)
}

// PreConsensusPTCWrongRootMsg carries a diverging observation (same block root, different
// payload-status flags): a valid message from a peer that did not converge, so it must fail the
// operator's expected-root check rather than reconstruct into a mixed signature.
var PreConsensusPTCWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return ptcMsg(msgSK, msgID, 1, true, TestingDutySlotGloas, false)
}

var ptcMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	msgCnt int,
	divergingData bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(slot)
	d, _ := beacon.DomainData(epoch, types.DomainPTCAttester)

	data := TestingPayloadAttestationData(slot)
	if divergingData {
		data = TestingDivergingPayloadAttestationData(slot)
	}
	signed, root, _ := signer.SignBeaconObject(data, d, sk.GetPublicKey().Serialize(), types.DomainPTCAttester)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(data, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainPTCAttester)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.PTCAttesterPartialSig,
		Slot:     slot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
		msgs.Messages = append(msgs.Messages, msg)
	}
	return &msgs
}
