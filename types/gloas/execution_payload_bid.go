package gloas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

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
	ParentBlockHash       phase0.Hash32              `ssz-index:"0"  ssz-size:"32"`
	ParentBlockRoot       phase0.Root                `ssz-index:"1"  ssz-size:"32"`
	BlockHash             phase0.Hash32              `ssz-index:"2"  ssz-size:"32"`
	PrevRandao            phase0.Root                `ssz-index:"3"  ssz-size:"32"`
	FeeRecipient          bellatrix.ExecutionAddress `ssz-index:"4"  ssz-size:"20"`
	GasLimit              uint64                     `ssz-index:"5"`
	BuilderIndex          BuilderIndex               `ssz-index:"6"`
	Slot                  phase0.Slot                `ssz-index:"7"`
	Value                 phase0.Gwei                `ssz-index:"8"`
	ExecutionPayment      phase0.Gwei                `ssz-index:"9"`
	BlobKZGCommitments    []deneb.KZGCommitment      `ssz-index:"10" ssz-type:"progressive-list"`
	ExecutionRequestsRoot phase0.Root                `ssz-index:"11"`
}

// SignedExecutionPayloadBid wraps an ExecutionPayloadBid with the builder's (or, for self-build, the
// proposer's) signature.
type SignedExecutionPayloadBid struct {
	Message   *ExecutionPayloadBid
	Signature phase0.BLSSignature `ssz-size:"96"`
}
