package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// BlindedExecutionPayloadEnvelope is the blinded form of the Gloas ExecutionPayloadEnvelope that
// the §6 envelope-signing duty signs (SIP #94 §6): the full `payload` is replaced by
// PayloadRoot = hash_tree_root(payload). In SSZ merkleization every field subtree commits to
// hash_tree_root(field), so substituting the payload with its root in the same field position
// preserves the container root — the blinded envelope hashes to the same root as the full one, and
// a BLS signature over the blinded signing root is valid for the full SignedExecutionPayloadEnvelope.
// (This holds provided both forms merkleize with the same container shape; see the EIP-7688
// ProgressiveContainer note in SIP #94 §6.) It is disseminated inside types.EnvelopeDissemination and
// threshold-signed under DomainBeaconBuilder — no QBFT consensus (SIP #94 §6). Blinding keeps that
// message bounded — a few hundred bytes rather than the full payload's hundreds of KB to ~MB. The full
// envelope/payload types are not vendored here (node-side only); the root-equivalence property is
// exercised by the node's tests against its full types.
type BlindedExecutionPayloadEnvelope struct {
	PayloadRoot phase0.Root `ssz-index:"0" ssz-size:"32"`
	// Gloas execution requests — the EIP-8282 five-list variant, not electra's three (see execution_requests.go).
	ExecutionRequests     *ExecutionRequests `ssz-index:"1"`
	BuilderIndex          BuilderIndex       `ssz-index:"2"`
	BeaconBlockRoot       phase0.Root        `ssz-index:"3" ssz-size:"32"`
	ParentBeaconBlockRoot phase0.Root        `ssz-index:"4" ssz-size:"32"`
}

// Encode/Decode wrap SSZ (de)serialization — the form disseminated in the §6 duty.
func (b *BlindedExecutionPayloadEnvelope) Encode() ([]byte, error)  { return b.MarshalSSZ() }
func (b *BlindedExecutionPayloadEnvelope) Decode(data []byte) error { return b.UnmarshalSSZ(data) }
