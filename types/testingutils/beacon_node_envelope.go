package testingutils

import (
	"errors"

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

// TestingBlindedExecutionPayloadEnvelope is the disseminated blinded envelope for the slot. Its bound
// fields match the fixture §4 block (gloas.TestingBeaconBlock): BeaconBlockRoot is the block's root,
// ParentBeaconBlockRoot is the block's parent root, and the empty ExecutionRequests hashes to the same
// root the block's bid commits — so the §6 binding checks pass.
var TestingBlindedExecutionPayloadEnvelope = func(slot phase0.Slot) *gloas.BlindedExecutionPayloadEnvelope {
	block := gloas.TestingBeaconBlock(slot)
	return &gloas.BlindedExecutionPayloadEnvelope{
		PayloadRoot:           phase0.Root{0x9a, 0x10, 0x0a, 0xd0},
		ExecutionRequests:     &gloas.ExecutionRequests{},
		BuilderIndex:          gloas.BuilderIndexSelfBuild,
		BeaconBlockRoot:       TestingProposedGloasBlockRoot(slot),
		ParentBeaconBlockRoot: block.ParentRoot,
	}
}

var TestingBlindedExecutionPayloadEnvelopeBytes = func(slot phase0.Slot) []byte {
	byts, err := TestingBlindedExecutionPayloadEnvelope(slot).Encode()
	if err != nil {
		panic(err.Error())
	}
	return byts
}

// TestingEnvelopeDissemination is the §6 dissemination carrier the builder operator broadcasts: the
// blinded envelope of the decided block, stamped with the duty slot.
var TestingEnvelopeDissemination = func(slot phase0.Slot) *types.EnvelopeDissemination {
	return &types.EnvelopeDissemination{
		Slot:     slot,
		Envelope: TestingBlindedExecutionPayloadEnvelope(slot),
	}
}

var TestingEnvelopeDisseminationBytes = func(slot phase0.Slot) []byte {
	byts, err := TestingEnvelopeDissemination(slot).Encode()
	if err != nil {
		panic(err.Error())
	}
	return byts
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

// TestingEnvelopeProposerNonBuilderDuty is an envelope duty for the slot at which this operator is not
// the builder operator (GetBlindedExecutionPayloadEnvelope errors), used for the no-publish path.
var TestingEnvelopeProposerNonBuilderDuty = func() *types.ValidatorDuty {
	duty := TestingEnvelopeProposerDuty()
	duty.Slot = TestingEnvelopeNonBuilderSlot
	return duty
}

// TestingEnvelopeNonBuilderSlot is the one slot for which the testing beacon node reports it did not
// build the decided block: GetBlindedExecutionPayloadEnvelope errors, modeling a non-builder operator
// that only signs a peer's dissemination and never publishes (SIP #94 §6). Keyed by slot so it survives
// the JSON test-vector round-trip, the same shape as TestingPTCAbstainSlot.
const TestingEnvelopeNonBuilderSlot = TestingDutySlotGloas + 3

// GetBlindedExecutionPayloadEnvelope returns this operator's produced blinded envelope for the slot's
// decided block. At TestingEnvelopeNonBuilderSlot the operator's beacon node did not build the block, so
// it returns an error — the runner then only signs a peer's dissemination (SIP #94 §6).
func (bn *TestingBeaconNode) GetBlindedExecutionPayloadEnvelope(slot phase0.Slot, blockRoot phase0.Root) (*gloas.BlindedExecutionPayloadEnvelope, error) {
	if slot == TestingEnvelopeNonBuilderSlot {
		return nil, errors.New("beacon node did not build the decided block")
	}
	envelope := TestingBlindedExecutionPayloadEnvelope(slot)
	envelope.BeaconBlockRoot = blockRoot
	return envelope, nil
}

// SubmitExecutionPayloadEnvelope records the published reveal's signing root (the blinded envelope's
// root). Publication happens only after the threshold signature reconstructs and verifies, so recording
// the envelope root is sufficient to confirm the builder operator published the selected envelope.
func (bn *TestingBeaconNode) SubmitExecutionPayloadEnvelope(envelope *gloas.BlindedExecutionPayloadEnvelope, signature phase0.BLSSignature) error {
	r, err := envelope.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}
