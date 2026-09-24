package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Regenerate with `go generate ./types/gloas/` (or `make generate-ssz`). --path . lets sszgen resolve the
// in-package Block field; --exclude-objs keeps it from regenerating the types other files own.
//go:generate sh -c "go run github.com/ferranbt/fastssz/sszgen --path . --include $(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/phase0,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/altair,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/capella,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/electra,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/bellatrix,$(go list -m -f '{{.Dir}}' github.com/attestantio/go-eth2-client)/spec/deneb --objs GloasProposalData --exclude-objs PayloadAttestation,BeaconBlockBody,BeaconBlock,SignedBeaconBlock,ExecutionPayloadBid,SignedExecutionPayloadBid,PayloadAttestationData,PayloadAttestationMessage,ExecutionRequests,BuilderDepositRequest,BuilderExitRequest,BlindedExecutionPayloadEnvelope,ProposerPreferences,SignedProposerPreferences,BuilderRequestAuth,SignedBuilderRequestAuth --output ./proposal_data_encoding.go"

// GloasProposalData is the SSVMessage.DataSSZ payload the SSV cluster decides at Gloas slots (SIP #94 §4):
// the Gloas beacon block plus the proposer's own self-build payload_root. payload_root is
// hash_tree_root(envelope.payload) — the one envelope field the block does not commit to — carried in the
// decided value so every operator derives and threshold-signs the §6 envelope root as a second entry of
// the block's post-consensus packet. Zero for an external bid (builder_index != BuilderIndexSelfBuild).
type GloasProposalData struct {
	Block       *BeaconBlock
	PayloadRoot phase0.Root `ssz-size:"32"`
}

// DecodeGloasProposalData unmarshals the Gloas DataSSZ wrapper (SIP #94 §4).
func DecodeGloasProposalData(dataSSZ []byte) (*GloasProposalData, error) {
	d := &GloasProposalData{}
	if err := d.UnmarshalSSZ(dataSSZ); err != nil {
		return nil, err
	}
	return d, nil
}
