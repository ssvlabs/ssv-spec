package gloas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Regenerate with `go generate ./types/gloas/` (or `make generate-ssz`). The includes are resolved from the module graph (`go list -m`)
// so they track go-eth2-client across dependency bumps rather than pinning.
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path ./execution_payload_bid.go --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/bellatrix,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/deneb --objs ExecutionPayloadBid,SignedExecutionPayloadBid"

// BuilderIndex identifies a builder in the Gloas builder registry. The sentinel BuilderIndexSelfBuild
// (BUILDER_INDEX_SELF_BUILD = UINT64_MAX) marks a self-built payload — the only path SSV produces.
type BuilderIndex uint64

// BuilderIndexSelfBuild (BUILDER_INDEX_SELF_BUILD) flags a self-built execution payload (SIP #94 §4).
const BuilderIndexSelfBuild = BuilderIndex(^uint64(0))

// MarshalJSON implements json.Marshaler, using the beacon-API decimal-string form. The self-build
// sentinel exceeds float64 precision, so a raw JSON number would not survive generic decoding.
func (b BuilderIndex) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatUint(uint64(b), 10))), nil
}

// UnmarshalJSON implements json.Unmarshaler, accepting the quoted beacon-API form and a bare number.
func (b *BuilderIndex) UnmarshalJSON(input []byte) error {
	unquoted := strings.Trim(string(input), `"`)
	index, err := strconv.ParseUint(unquoted, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid value for builder index: %w", err)
	}
	*b = BuilderIndex(index)
	return nil
}

// ExecutionPayloadBid is the Gloas (ePBS) bid the proposer commits to in the block body, replacing the
// inline execution payload of pre-Gloas blocks: the block carries only this commitment, and the payload
// itself ships separately in the envelope (§6). For self-build, BuilderIndex is BuilderIndexSelfBuild.
// Fields match the pinned spec / go-eth2-client PR #280.
type ExecutionPayloadBid struct {
	ParentBlockHash       phase0.Hash32              `ssz-size:"32"`
	ParentBlockRoot       phase0.Root                `ssz-size:"32"`
	BlockHash             phase0.Hash32              `ssz-size:"32"`
	PrevRandao            phase0.Hash32              `ssz-size:"32"`
	FeeRecipient          bellatrix.ExecutionAddress `ssz-size:"20"`
	GasLimit              uint64
	BuilderIndex          BuilderIndex
	Slot                  phase0.Slot
	Value                 phase0.Gwei
	ExecutionPayment      phase0.Gwei
	BlobKZGCommitments    []deneb.KZGCommitment `ssz-max:"4096" ssz-size:"?,48"`
	ExecutionRequestsRoot phase0.Root           `ssz-size:"32"`
}

// SignedExecutionPayloadBid wraps an ExecutionPayloadBid with the builder's (or, for self-build, the
// proposer's) signature.
type SignedExecutionPayloadBid struct {
	Message   *ExecutionPayloadBid
	Signature phase0.BLSSignature `ssz-size:"96"`
}
