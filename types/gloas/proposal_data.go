package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

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
