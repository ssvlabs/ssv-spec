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
		Epoch: TestingDutyEpoch,
		Root:  TestingBlockRoot,
	},
	Target: &phase0.Checkpoint{
		Epoch: TestingDutyEpoch + 1,
		Root:  TestingBlockRoot,
	},
}
var TestBeaconVoteByts, _ = TestBeaconVote.Encode()

var TestSlashableBeaconVote = types.BeaconVote{
	BlockRoot: phase0.Root{1, 2, 3, 4},
	Source: &phase0.Checkpoint{
		Epoch: TestingDutyEpoch + 1,
		Root:  phase0.Root{1, 2, 3, 4},
	},
	Target: &phase0.Checkpoint{
		Epoch: TestingDutyEpoch,
		Root:  phase0.Root{1, 2, 3, 5},
	},
}
var TestSlashableBeaconVoteByts, _ = TestSlashableBeaconVote.Encode()
