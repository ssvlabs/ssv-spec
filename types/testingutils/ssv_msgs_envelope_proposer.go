package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgEnvelopeProposer = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleEnvelopeProposer))
}

// ==================================================
// Consensus
// ==================================================

// SSVDecidingMsgsForEnvelopeProposer returns the §6 envelope QBFT deciding messages for the slot; the
// envelope duty has no pre-consensus phase.
var SSVDecidingMsgsForEnvelopeProposer = func(slot phase0.Slot, ks *TestKeySet) []*types.SignedSSVMessage {
	id := types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleEnvelopeProposer)
	fullData := TestingEnvelopeConsensusDataByts(slot)
	r, err := qbft.HashDataRoot(fullData)
	if err != nil {
		panic(err)
	}
	return SSVDecidingMsgsForHeightWithRoot(r, fullData, id[:], qbft.Height(slot), ks)
}

// ==================================================
// PostConsensus
// ==================================================

var PostConsensusEnvelopeProposerMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusEnvelopeMsg(sk, id, 1, false, false)
}

var PostConsensusEnvelopeProposerTooManyRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusEnvelopeMsg(sk, id, 2, false, false)
}

var PostConsensusEnvelopeProposerTooFewRootsMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusEnvelopeMsg(sk, id, 0, false, false)
}

var PostConsensusEnvelopeProposerWrongRootMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusEnvelopeMsg(sk, id, 1, true, false)
}

var PostConsensusEnvelopeProposerWrongBeaconSigMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return postConsensusEnvelopeMsg(sk, id, 1, false, true)
}

// PostConsensusEnvelopeProposerMsgForEnvelope signs the given decided envelope's root — used when the
// decided envelope is not this operator's own production.
var PostConsensusEnvelopeProposerMsgForEnvelope = func(sk *bls.SecretKey, id types.OperatorID, envelope *gloas.BlindedExecutionPayloadEnvelope) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()

	root, err := envelope.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	hashRoot := types.SSZ32Bytes(root)

	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(TestingDutySlotGloas)
	d, _ := beacon.DomainData(epoch, types.DomainBeaconBuilder)
	sig, signingRoot, _ := signer.SignBeaconObject(hashRoot, d, sk.GetPublicKey().Serialize(), types.DomainBeaconBuilder)
	blsSig := phase0.BLSSignature{}
	copy(blsSig[:], sig)

	return &types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingDutySlotGloas,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: blsSig[:],
				SigningRoot:      signingRoot,
				Signer:           id,
				ValidatorIndex:   TestingValidatorIndex,
			},
		},
	}
}

var postConsensusEnvelopeMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()

	envelope := TestingBlindedExecutionPayloadEnvelope(TestingDutySlotGloas)
	if wrongRoot {
		envelope.PayloadRoot = phase0.Root{0xba, 0xdb, 0xad}
	}
	root, err := envelope.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	hashRoot := types.SSZ32Bytes(root)

	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(TestingDutySlotGloas)
	d, _ := beacon.DomainData(epoch, types.DomainBeaconBuilder)
	sig, signingRoot, _ := signer.SignBeaconObject(hashRoot, d, sk.GetPublicKey().Serialize(), types.DomainBeaconBuilder)
	if wrongBeaconSig {
		sig, signingRoot, _ = signer.SignBeaconObject(hashRoot, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainBeaconBuilder)
	}
	blsSig := phase0.BLSSignature{}
	copy(blsSig[:], sig)

	msgs := types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingDutySlotGloas,
		Messages: []*types.PartialSignatureMessage{},
	}
	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: blsSig[:],
			SigningRoot:      signingRoot,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
		msgs.Messages = append(msgs.Messages, msg)
	}
	return &msgs
}
