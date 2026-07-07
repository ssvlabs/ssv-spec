package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Regenerate with `go generate ./types/gloas/` (or `make generate-ssz`). --path is the package dir
// (not just this file) so sszgen resolves the sibling gloas types the envelope references; --objs
// limits output to the blinded envelope forms, collected into their own _encoding.go. Includes
// track go-eth2-client via `go list -m`.
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path . --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/electra,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/bellatrix,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/capella --objs BlindedExecutionPayloadEnvelope,SignedBlindedExecutionPayloadEnvelope --exclude-objs ExecutionRequests,BuilderDepositRequest,BuilderExitRequest --output ./execution_payload_envelope_encoding.go"

// BlindedExecutionPayloadEnvelope is the blinded form of the Gloas ExecutionPayloadEnvelope that
// the §6 envelope-signing duty signs (SIP #94 §6): the full `payload` is replaced by
// PayloadRoot = hash_tree_root(payload). In SSZ merkleization every field subtree commits to
// hash_tree_root(field), so substituting the payload with its root in the same field position
// preserves the container root — the blinded envelope hashes to the same root as the full one, and
// a BLS signature over the blinded signing root is valid for the full SignedExecutionPayloadEnvelope.
// (This holds provided both forms merkleize with the same container shape; see the EIP-7688
// ProgressiveContainer note in SIP #94 §6.) It rides in EnvelopeConsensusData.DataSSZ; blinding
// keeps that QBFT value bounded — a few hundred bytes rather than the full payload's hundreds of
// KB to ~MB. The full envelope/payload types are not vendored here (node-side only); the
// root-equivalence property is exercised by the node's tests against its full types.
type BlindedExecutionPayloadEnvelope struct {
	PayloadRoot phase0.Root `ssz-size:"32"`
	// Gloas execution requests — the EIP-8282 five-list variant, not electra's three (see execution_requests.go).
	ExecutionRequests     *ExecutionRequests
	BuilderIndex          BuilderIndex
	BeaconBlockRoot       phase0.Root `ssz-size:"32"`
	ParentBeaconBlockRoot phase0.Root `ssz-size:"32"`
}

// Encode/Decode wrap SSZ (de)serialization — the form carried in the §6 QBFT consensus DataSSZ.
func (b *BlindedExecutionPayloadEnvelope) Encode() ([]byte, error)  { return b.MarshalSSZ() }
func (b *BlindedExecutionPayloadEnvelope) Decode(data []byte) error { return b.UnmarshalSSZ(data) }

// SignedBlindedExecutionPayloadEnvelope wraps the blinded envelope with the builder's signature
// (under DOMAIN_BEACON_BUILDER) — the §6 publication body on the blinded (stateful) path, where the
// producing beacon node reconstructs the full envelope from its cache. The signature is valid here
// because the blinded root equals the full envelope's (see BlindedExecutionPayloadEnvelope) — the
// same property the §6 duty relies on to sign the blinded form.
type SignedBlindedExecutionPayloadEnvelope struct {
	Message   *BlindedExecutionPayloadEnvelope
	Signature phase0.BLSSignature `ssz-size:"96"`
}
