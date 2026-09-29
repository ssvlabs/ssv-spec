package gloas

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSignedBeaconBlockDevnet6Fields cross-checks the decoded devnet-6 block against facts known independently of
// this codec. Its byte round-trip is a spec vector, but that uses the same schema both ways, so it can't catch two
// same-size fields swapped inside the bid; these decoded-field checks can.
func TestSignedBeaconBlockDevnet6Fields(t *testing.T) {
	var blk SignedBeaconBlock
	require.NoError(t, blk.UnmarshalSSZ(TestingDevnet6SignedBeaconBlockSSZ), "decode real v8.2.0 Gloas block")
	require.NotNil(t, blk.Message.Body.ParentExecutionRequests, "Gloas block body carries execution requests")
	require.NotNil(t, blk.Message.Body.SignedExecutionPayloadBid, "Gloas block body carries the execution-payload bid")

	// The known slot, and the bid commitments that must equal the block's own.
	bid := blk.Message.Body.SignedExecutionPayloadBid.Message
	require.EqualValues(t, 66, blk.Message.Slot, "devnet-6 fixture is slot 66")
	require.Equal(t, blk.Message.Slot, bid.Slot, "bid commits to the block's slot")
	require.Equal(t, blk.Message.ParentRoot, bid.ParentBlockRoot, "bid commits to the block's parent root")
}
