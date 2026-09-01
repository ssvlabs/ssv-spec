package gloas

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// TestSignedBeaconBlockRootMatchesGoEth2Client pins the core invariant of the dynamic-ssz migration:
// ssv-spec's Gloas SignedBeaconBlock must produce the same hash_tree_root as go-eth2-client's
// spec/gloas type — the progressive-merkleization reference validated against consensus-specs on
// devnet-8. Decoding the same real block into both and comparing roots catches a future dynssz/tag
// drift that the round-trip and wire tests would miss: serialization can byte-match while HTR
// diverges, which is exactly the class of bug (wrong progressive root) this migration fixes.
//
// This exercises the progressive *container* merkleization (the block body). The devnet block's
// progressive *lists* are all empty, so TestExecutionPayloadBidRootMatchesGoEth2Client below covers
// a non-empty progressive list, where progressive and binary merkleization actually diverge.
func TestSignedBeaconBlockRootMatchesGoEth2Client(t *testing.T) {
	var ours SignedBeaconBlock
	require.NoError(t, ours.UnmarshalSSZ(devnet6GloasBlockSSZ), "decode into ssv-spec gloas type")
	ourRoot, err := ours.HashTreeRoot()
	require.NoError(t, err)

	var ref eth2gloas.SignedBeaconBlock
	require.NoError(t, ref.UnmarshalSSZ(devnet6GloasBlockSSZ), "decode into go-eth2-client spec/gloas type")
	refRoot, err := ref.HashTreeRoot()
	require.NoError(t, err)

	require.Equal(t, refRoot, ourRoot,
		"ssv-spec Gloas SignedBeaconBlock hash_tree_root must match go-eth2-client spec/gloas (progressive merkleization)")
}

// TestExecutionPayloadBidRootMatchesGoEth2Client guards progressive-*list* merkleization: with a
// non-empty BlobKZGCommitments (a `ssz-type:"progressive-list"`), a progressive vs a binary
// MerkleizeWithMixin produce different roots, so this fails loudly if the tag/backend ever drifts.
// A real block with empty lists (the test above) cannot catch that.
func TestExecutionPayloadBidRootMatchesGoEth2Client(t *testing.T) {
	bid := &ExecutionPayloadBid{
		GasLimit:           30_000_000,
		BuilderIndex:       BuilderIndexSelfBuild,
		Slot:               66,
		Value:              12345,
		ExecutionPayment:   678,
		BlobKZGCommitments: []deneb.KZGCommitment{{0x01}, {0x02}, {0x03}, {0x04}, {0x05}},
	}
	ssz, err := bid.MarshalSSZ()
	require.NoError(t, err)

	var ours ExecutionPayloadBid
	require.NoError(t, ours.UnmarshalSSZ(ssz), "decode into ssv-spec gloas type")
	ourRoot, err := ours.HashTreeRoot()
	require.NoError(t, err)

	var ref eth2gloas.ExecutionPayloadBid
	require.NoError(t, ref.UnmarshalSSZ(ssz), "decode into go-eth2-client spec/gloas type")
	refRoot, err := ref.HashTreeRoot()
	require.NoError(t, err)

	require.Equal(t, refRoot, ourRoot,
		"non-empty progressive-list root must match go-eth2-client spec/gloas")
}
