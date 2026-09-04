package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgEnvelopeProposer = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleEnvelopeProposer))
}

// ==================================================
// Dissemination (SIP #94 §6)
// ==================================================

// EnvelopeDisseminationSSVMsg is the §6 dissemination SSVMessage carrying the slot's blinded envelope.
var EnvelopeDisseminationSSVMsg = func(slot phase0.Slot) *types.SSVMessage {
	return &types.SSVMessage{
		MsgType: types.SSVEnvelopeDisseminationMsgType,
		MsgID:   types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleEnvelopeProposer),
		Data:    TestingEnvelopeDisseminationBytes(slot),
	}
}

// SignedEnvelopeDisseminationSSVMsg is a §6 dissemination operator-signed by the given operator (the
// builder operator disseminating its blinded envelope).
var SignedEnvelopeDisseminationSSVMsg = func(ks *TestKeySet, signer types.OperatorID, slot phase0.Slot) *types.SignedSSVMessage {
	return SignedSSVMessageWithSigner(signer, ks.OperatorKeys[signer], EnvelopeDisseminationSSVMsg(slot))
}

// ==================================================
// PreConsensus (the single threshold-signing round, EnvelopePartialSig)
// ==================================================

var PreConsensusEnvelopeMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 1, false, TestingDutySlotGloas, false)
}

var PreConsensusEnvelopeTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 0, false, TestingDutySlotGloas, false)
}

var PreConsensusEnvelopeTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 2, false, TestingDutySlotGloas, false)
}

var PreConsensusEnvelopeWrongBeaconSigMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 1, false, TestingDutySlotGloas, true)
}

// PreConsensusEnvelopeWrongRootMsg signs a diverging envelope (different PayloadRoot): a peer that
// selected a different envelope, whose root fails the operator's expected-root check rather than
// reconstructing into a mixed signature.
var PreConsensusEnvelopeWrongRootMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 1, true, TestingDutySlotGloas, false)
}

var PreConsensusEnvelopeNextEpochMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 1, false, TestingDutySlotGloasNextEpoch, false)
}

// PreConsensusEnvelopeMsgForSlot signs the slot's blinded envelope root — used at the non-builder slot,
// where the operator signs a peer's disseminated envelope but does not publish.
var PreConsensusEnvelopeMsgForSlot = func(sk *bls.SecretKey, id types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	return envelopeMsg(sk, id, 1, false, slot, false)
}

var envelopeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	msgCnt int,
	divergingEnvelope bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(slot)
	d, _ := beacon.DomainData(epoch, types.DomainBeaconBuilder)

	envelope := TestingBlindedExecutionPayloadEnvelope(slot)
	if divergingEnvelope {
		envelope.PayloadRoot = phase0.Root{0xba, 0xdb, 0xad}
	}

	signed, root, _ := signer.SignBeaconObject(envelope, d, sk.GetPublicKey().Serialize(), types.DomainBeaconBuilder)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(envelope, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainBeaconBuilder)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.EnvelopePartialSig,
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
