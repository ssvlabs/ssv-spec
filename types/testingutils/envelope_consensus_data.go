package testingutils

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// EnvelopeConsensusData (ePBS / SIP #94 §6)
// ==================================================

// TestBlindedExecutionPayloadEnvelope is a minimal self-build blinded envelope used as the §6 QBFT
// value's DataSSZ. SSZ serialization is stable regardless of the EIP-7688 progressive-merkleization
// question; that question only affects the blinded envelope's own structured hash tree root (used in
// §6 signing, PR 7), not its serialized bytes.
var TestBlindedExecutionPayloadEnvelope = &gloas.BlindedExecutionPayloadEnvelope{
	PayloadRoot:           TestingBlockRoot,
	ExecutionRequests:     &gloas.ExecutionRequests{},
	BuilderIndex:          gloas.BuilderIndexSelfBuild,
	BeaconBlockRoot:       TestingBlockRoot,
	ParentBeaconBlockRoot: TestingBlockRoot,
}

var testEnvelopeDataSSZ, _ = TestBlindedExecutionPayloadEnvelope.Encode()

// TestEnvelopeConsensusData is the §6 envelope-signing QBFT value. EnvelopeConsensusData is a plain
// container whose DataSSZ is opaque bytes, so its hash tree root is stable (non-provisional) even
// though DataSSZ carries a (progressive) blinded envelope.
var TestEnvelopeConsensusData = types.EnvelopeConsensusData{
	Duty: types.ValidatorDuty{
		Type:           types.BNRoleEnvelopeProposer,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlot,
		ValidatorIndex: TestingValidatorIndex,
	},
	Version: gloas.DataVersionGloas,
	DataSSZ: testEnvelopeDataSSZ,
}
