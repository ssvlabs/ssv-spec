package gloas

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// TestingBeaconBlock returns a minimal self-build Gloas BeaconBlock for the slot, with the required
// fixed-size body fields populated so it round-trips through SSZ. The bid commits the execution-requests
// root of an empty ExecutionRequests; the §6 blinded envelope derived from the decided block copies that
// same root straight from the bid (SIP #94 §6), so its post-consensus signing root is well-defined. For
// use in tests.
func TestingBeaconBlock(slot phase0.Slot) *BeaconBlock {
	requestsRoot, err := (&ExecutionRequests{}).HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}
	return &BeaconBlock{
		Slot: slot,
		Body: &BeaconBlockBody{
			ETH1Data:      &phase0.ETH1Data{BlockHash: make([]byte, 32)},
			SyncAggregate: &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &SignedExecutionPayloadBid{Message: &ExecutionPayloadBid{
				BuilderIndex:          BuilderIndexSelfBuild,
				ExecutionRequestsRoot: requestsRoot,
			}},
			ParentExecutionRequests: &ExecutionRequests{},
		},
	}
}

// TestingExternalBuilderIndex is a non-self-build BuilderIndex for the §4 external-bid path — any value
// other than BuilderIndexSelfBuild marks the bid as won by an external builder.
const TestingExternalBuilderIndex = BuilderIndex(3)

// TestingBeaconBlockExternalBuild returns the minimal Gloas block of TestingBeaconBlock but with an
// external-builder bid (BuilderIndex != BuilderIndexSelfBuild): the §4 path where produce yields a bare
// block and no envelope, so the decided value's payload_root is zero. For use in tests.
func TestingBeaconBlockExternalBuild(slot phase0.Slot) *BeaconBlock {
	blk := TestingBeaconBlock(slot)
	blk.Body.SignedExecutionPayloadBid.Message.BuilderIndex = TestingExternalBuilderIndex
	return blk
}
