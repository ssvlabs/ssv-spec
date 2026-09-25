package gloas

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

// devnet6GloasBlockSSZ is a real on-chain Gloas SignedBeaconBlock (slot 66) captured from lighthouse
// v8.2.0 (ethpandaops glamsterdam-devnet-6). Its ParentExecutionRequests carries the full EIP-8282
// five-list ExecutionRequests (all lists empty on this block).
//
//go:embed testdata/devnet6_gloas_block.ssz
var devnet6GloasBlockSSZ []byte

// TestSignedBeaconBlockMatchesDevnet6Wire guards that the SignedBeaconBlock codec byte-round-trips
// a real Glamsterdam CL's wire format — the check that pins the §4 submit. A future wire drift (a new
// request list, a reordered field) fails here instead of only in a devnet run. It also asserts the slot
// and the bid's consensus invariants against the block: the byte round-trip uses the same schema both
// ways, so it can't catch two same-size fields swapped inside the bid — these decoded-field checks can.
func TestSignedBeaconBlockMatchesDevnet6Wire(t *testing.T) {
	var blk SignedBeaconBlock
	require.NoError(t, blk.UnmarshalSSZ(devnet6GloasBlockSSZ), "decode real v8.2.0 Gloas block")
	require.NotNil(t, blk.Message.Body.ParentExecutionRequests, "Gloas block body carries execution requests")
	require.NotNil(t, blk.Message.Body.SignedExecutionPayloadBid, "Gloas block body carries the execution-payload bid")

	// Decoded-field checks the byte round-trip can't make: the known slot, and the bid commitments that
	// must equal the block's own.
	bid := blk.Message.Body.SignedExecutionPayloadBid.Message
	require.EqualValues(t, 66, blk.Message.Slot, "devnet-6 fixture is slot 66")
	require.Equal(t, blk.Message.Slot, bid.Slot, "bid commits to the block's slot")
	require.Equal(t, blk.Message.ParentRoot, bid.ParentBlockRoot, "bid commits to the block's parent root")

	out, err := blk.MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, len(devnet6GloasBlockSSZ), len(out), "re-marshal length must match the CL wire size")
	require.Equal(t, devnet6GloasBlockSSZ, out,
		"re-marshal must byte-match v8.2.0's wire format — a mismatch means the Gloas types drifted from the CL")
}
