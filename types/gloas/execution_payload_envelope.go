package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Regenerate with `go generate ./types/gloas/` (or `make generate-ssz`). --path is the package dir
// (not just this file) so sszgen resolves the sibling gloas types the envelope references; --objs
// limits output to the blinded envelope forms, collected into their own _encoding.go. Includes
// track go-eth2-client via `go list -m`.
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path . --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/electra,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/bellatrix,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/capella --objs BlindedExecutionPayloadEnvelope --exclude-objs ExecutionRequests,BuilderDepositRequest,BuilderExitRequest --output ./execution_payload_envelope_encoding.go"

// BlindedExecutionPayloadEnvelope is the blinded form of the Gloas ExecutionPayloadEnvelope that the §6
// envelope signing signs (SIP #94 §6): the full `payload` and `execution_requests` are replaced by their
// roots (PayloadRoot, ExecutionRequestsRoot). In SSZ merkleization every field subtree commits to
// hash_tree_root(field), so substituting each field with its root in the same position preserves the
// container root — the blinded envelope hashes to the same root as the full one, and a BLS signature over
// the blinded signing root is valid for the full SignedExecutionPayloadEnvelope. (This holds provided both
// forms merkleize with the same container shape; see the EIP-7688 ProgressiveContainer note in SIP #94 §6.)
// Every operator derives it from the §4-decided value — payload_root from the decided value,
// execution_requests_root from the bid — and threshold-signs its root under DomainBeaconBuilder as a second
// entry of the block's post-consensus packet, no separate round. Blinding keeps the payload and the
// unbounded requests list off the SSV wire. The full envelope/payload types are not vendored here
// (node-side only); the root-equivalence property is exercised by the node's tests against its full types.
type BlindedExecutionPayloadEnvelope struct {
	PayloadRoot           phase0.Root `ssz-size:"32"`
	ExecutionRequestsRoot phase0.Root `ssz-size:"32"`
	BuilderIndex          BuilderIndex
	BeaconBlockRoot       phase0.Root `ssz-size:"32"`
	ParentBeaconBlockRoot phase0.Root `ssz-size:"32"`
}
