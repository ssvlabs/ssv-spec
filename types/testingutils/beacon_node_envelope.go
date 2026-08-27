package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// Execution Payload Envelope (SIP #94 §6)
// ==================================================

// TestingProposedGloasBlockRoot is the §4-decided block root for the slot — the root of the fixture
// block the proposer vectors decide on, tying the envelope fixtures to the same block.
var TestingProposedGloasBlockRoot = func(slot phase0.Slot) phase0.Root {
	root, err := gloas.TestingBeaconBlock(slot).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}
	return root
}

var TestingBlindedExecutionPayloadEnvelope = func(slot phase0.Slot) *gloas.BlindedExecutionPayloadEnvelope {
	return &gloas.BlindedExecutionPayloadEnvelope{
		PayloadRoot:           phase0.Root{0x9a, 0x10, 0x0a, 0xd0},
		ExecutionRequests:     &gloas.ExecutionRequests{},
		BuilderIndex:          gloas.BuilderIndexSelfBuild,
		BeaconBlockRoot:       TestingProposedGloasBlockRoot(slot),
		ParentBeaconBlockRoot: TestingBlockRoot,
	}
}

var TestingBlindedExecutionPayloadEnvelopeBytes = func(slot phase0.Slot) []byte {
	byts, err := TestingBlindedExecutionPayloadEnvelope(slot).Encode()
	if err != nil {
		panic(err.Error())
	}
	return byts
}

// TestingSignedBlindedExecutionPayloadEnvelope is the published object: the blinded envelope with the
// reconstructed validator signature under DomainBeaconBuilder.
var TestingSignedBlindedExecutionPayloadEnvelope = func(ks *TestKeySet, slot phase0.Slot) *gloas.SignedBlindedExecutionPayloadEnvelope {
	envelope := TestingBlindedExecutionPayloadEnvelope(slot)
	return &gloas.SignedBlindedExecutionPayloadEnvelope{
		Message:   envelope,
		Signature: signBeaconObject(envelope, types.DomainBeaconBuilder, ks),
	}
}

var TestingEnvelopeProposerDuty = func() *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:           types.BNRoleEnvelopeProposer,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlotGloas,
		ValidatorIndex: TestingValidatorIndex,
	}
}

var TestingEnvelopeProposerNextEpochDuty = func() *types.ValidatorDuty {
	duty := TestingEnvelopeProposerDuty()
	duty.Slot = TestingDutySlotGloasNextEpoch
	return duty
}

var TestingEnvelopeConsensusData = func(slot phase0.Slot) *types.EnvelopeConsensusData {
	duty := TestingEnvelopeProposerDuty()
	duty.Slot = slot
	return &types.EnvelopeConsensusData{
		Duty:    *duty,
		Version: gloas.DataVersionGloas,
		DataSSZ: TestingBlindedExecutionPayloadEnvelopeBytes(slot),
	}
}

var TestingEnvelopeConsensusDataByts = func(slot phase0.Slot) []byte {
	byts, err := TestingEnvelopeConsensusData(slot).Encode()
	if err != nil {
		panic(err.Error())
	}
	return byts
}

// GetBlindedExecutionPayloadEnvelope returns this operator's produced blinded envelope for the slot's
// decided block
func (bn *TestingBeaconNode) GetBlindedExecutionPayloadEnvelope(slot phase0.Slot, blockRoot phase0.Root) (*gloas.BlindedExecutionPayloadEnvelope, error) {
	envelope := TestingBlindedExecutionPayloadEnvelope(slot)
	envelope.BeaconBlockRoot = blockRoot
	return envelope, nil
}

// SubmitBlindedExecutionPayloadEnvelope records the signed envelope's root
func (bn *TestingBeaconNode) SubmitBlindedExecutionPayloadEnvelope(envelope *gloas.SignedBlindedExecutionPayloadEnvelope) error {
	r, err := envelope.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}
