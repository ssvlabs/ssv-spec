package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// EnvelopeDissemination is the SSVMessage.Data payload for SSVEnvelopeDisseminationMsgType (SIP #94 §6):
// the Gloas (ePBS) self-build builder operator broadcasts the blinded execution-payload envelope of the
// §4-decided block for every operator to content-select and threshold-sign. Slot stamps the duty slot
// for message validation, mirroring PartialSignatureMessages.Slot. Envelope is the blinded form (payload
// replaced by its root), whose root equals the full envelope's, so a signature over it is valid for the
// full SignedExecutionPayloadEnvelope.
type EnvelopeDissemination struct {
	Slot     phase0.Slot
	Envelope *gloas.BlindedExecutionPayloadEnvelope
}

// Encode the EnvelopeDissemination object
func (d *EnvelopeDissemination) Encode() ([]byte, error) {
	return d.MarshalSSZ()
}

// Decode the EnvelopeDissemination object
func (d *EnvelopeDissemination) Decode(data []byte) error {
	return d.UnmarshalSSZ(data)
}
