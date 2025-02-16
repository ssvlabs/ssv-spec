package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// BeaconVote
// ==================================================

var TestBeaconVote = types.BeaconVote{
	BlockRoot: TestingBlockRoot,
	Source: &phase0.Checkpoint{
		Epoch: 0,
		Root:  TestingBlockRoot,
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  TestingBlockRoot,
	},
}
var TestBeaconVoteByts, _ = TestBeaconVote.Encode()

var TestBeaconVoteNextEpoch = types.BeaconVote{
	BlockRoot: TestingBlockRoot,
	Source: &phase0.Checkpoint{
		Epoch: 0,
		Root:  TestingBlockRoot,
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  TestingBlockRoot,
	},
}
var TestBeaconVoteNextEpochByts, _ = TestBeaconVoteNextEpoch.Encode()

var TestWrongBeaconVote = types.BeaconVote{
	BlockRoot: phase0.Root{1, 2, 3, 4},
	Source: &phase0.Checkpoint{
		Epoch: 2,
		Root:  phase0.Root{1, 2, 3, 4},
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  phase0.Root{1, 2, 3, 5},
	},
}
var TestWrongBeaconVoteByts, _ = TestWrongBeaconVote.Encode()
