package testingutils

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// EnvelopeDissemination (ePBS / SIP #94 §6)
// ==================================================

// TestBlindedExecutionPayloadEnvelope is a minimal self-build blinded envelope carried by the §6
// dissemination message. Its SSZ serialization is stable regardless of the EIP-7688 progressive-
// merkleization question; that question only affects the blinded envelope's own structured hash tree
// root (used in §6 signing), not the disseminated carrier's serialized bytes.
var TestBlindedExecutionPayloadEnvelope = &gloas.BlindedExecutionPayloadEnvelope{
	PayloadRoot:           TestingBlockRoot,
	ExecutionRequests:     &gloas.ExecutionRequests{},
	BuilderIndex:          gloas.BuilderIndexSelfBuild,
	BeaconBlockRoot:       TestingBlockRoot,
	ParentBeaconBlockRoot: TestingBlockRoot,
}

// TestEnvelopeDissemination is the §6 dissemination carrier (SSVMessage.Data payload for
// SSVEnvelopeDisseminationMsgType) used by the encoding suite.
var TestEnvelopeDissemination = types.EnvelopeDissemination{
	Slot:     TestingDutySlot,
	Envelope: TestBlindedExecutionPayloadEnvelope,
}
