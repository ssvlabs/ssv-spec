package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// PTC — Payload Timeliness Committee (SIP #94 §3)
// ==================================================

// TestingPTCAbstainSlot is the one slot for which the testing beacon node reports no block (a zero
// BeaconBlockRoot) — the §3 abstain trigger. Every other slot yields TestingPayloadAttestationData.
const TestingPTCAbstainSlot = TestingDutySlotGloas + 2

var TestingPayloadAttestationData = func(slot phase0.Slot) *gloas.PayloadAttestationData {
	return &gloas.PayloadAttestationData{
		BeaconBlockRoot:   TestingBlockRoot,
		Slot:              slot,
		PayloadPresent:    true,
		BlobDataAvailable: true,
	}
}

// TestingDivergingPayloadAttestationData is a valid observation that differs from
// TestingPayloadAttestationData only in the payload-status flags — a peer that saw the block but not
// the payload. Used to model divergence in the honest-convergence checks.
var TestingDivergingPayloadAttestationData = func(slot phase0.Slot) *gloas.PayloadAttestationData {
	return &gloas.PayloadAttestationData{
		BeaconBlockRoot: TestingBlockRoot,
		Slot:            slot,
	}
}

var TestingPTCAttesterDuty = func() *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:           types.BNRolePTCAttester,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlotGloas,
		ValidatorIndex: TestingValidatorIndex,
	}
}

var TestingPTCAttesterAbstainDuty = func() *types.ValidatorDuty {
	duty := TestingPTCAttesterDuty()
	duty.Slot = TestingPTCAbstainSlot
	return duty
}

var TestingPTCAttesterNextEpochDuty = func() *types.ValidatorDuty {
	duty := TestingPTCAttesterDuty()
	duty.Slot = TestingDutySlotGloasNextEpoch
	return duty
}

// TestingSignedPayloadAttestationMessage is the reconstructed message the runner submits on quorum.
var TestingSignedPayloadAttestationMessage = func(ks *TestKeySet) *gloas.PayloadAttestationMessage {
	data := TestingPayloadAttestationData(TestingDutySlotGloas)
	return &gloas.PayloadAttestationMessage{
		ValidatorIndex: TestingValidatorIndex,
		Data:           data,
		Signature:      signBeaconObject(data, types.DomainPTCAttester, ks),
	}
}

// GetPayloadAttestationData returns the slot's payload attestation data; TestingPTCAbstainSlot
// reports a zero BeaconBlockRoot — the "no block seen" abstain contract (SIP #94 §3).
func (bn *TestingBeaconNode) GetPayloadAttestationData(slot phase0.Slot) (*gloas.PayloadAttestationData, error) {
	if slot == TestingPTCAbstainSlot {
		return &gloas.PayloadAttestationData{Slot: slot}, nil
	}
	return TestingPayloadAttestationData(slot), nil
}

// SubmitPayloadAttestation records the payload attestation message's root
func (bn *TestingBeaconNode) SubmitPayloadAttestation(msg *gloas.PayloadAttestationMessage) error {
	r, err := msg.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}
