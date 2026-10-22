package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// BlindedExecutionPayloadEnvelope is the blinded form of the Gloas ExecutionPayloadEnvelope that the §6
// envelope signing signs (SIP #94 §6): the full `payload` and `execution_requests` are replaced by their
// roots (PayloadRoot, ExecutionRequestsRoot). In SSZ merkleization every field subtree commits to
// hash_tree_root(field), so substituting each field with its root in the same position preserves the
// container root — the blinded envelope hashes to the same root as the full one, and a BLS signature over
// the blinded signing root is valid for the full SignedExecutionPayloadEnvelope. (This holds provided both
// forms merkleize with the same progressive-container shape — the ssz-index tags below; see the EIP-7688
// ProgressiveContainer note in SIP #94 §6.) Every operator derives it from the §4-decided value —
// payload_root from the decided value, execution_requests_root from the bid — and threshold-signs its root
// under DomainBeaconBuilder as a second entry of the block's post-consensus packet, no separate round.
// Blinding keeps the payload and the unbounded requests list off the SSV wire. The full envelope/payload
// types are not vendored here (node-side only); the root-equivalence property is exercised by the node's
// tests against its full types. The dynssz-generated encoder lives in blindedexecutionpayloadenvelope_ssz.go.
type BlindedExecutionPayloadEnvelope struct {
	PayloadRoot           phase0.Root  `ssz-index:"0" ssz-size:"32"`
	ExecutionRequestsRoot phase0.Root  `ssz-index:"1" ssz-size:"32"`
	BuilderIndex          BuilderIndex `ssz-index:"2"`
	BeaconBlockRoot       phase0.Root  `ssz-index:"3" ssz-size:"32"`
	ParentBeaconBlockRoot phase0.Root  `ssz-index:"4" ssz-size:"32"`
}
