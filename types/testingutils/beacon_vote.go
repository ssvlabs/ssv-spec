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

// ==================================================
// GloasBeaconVote (ePBS / SIP #94 §2)
// ==================================================

// TestGloasBeaconVote mirrors TestBeaconVote and adds the BN-supplied AttestationDataIndex
// (1 = payload FULL). Fixed 120-byte SSZ; plain container, so its hash tree root is stable
// (unaffected by the EIP-7688 progressive-merkleization question, which only touches §4/§6).
var TestGloasBeaconVote = types.GloasBeaconVote{
	BlockRoot: TestingBlockRoot,
	Source: &phase0.Checkpoint{
		Epoch: 0,
		Root:  TestingBlockRoot,
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  TestingBlockRoot,
	},
	AttestationDataIndex: 1,
}
var TestGloasBeaconVoteByts, _ = TestGloasBeaconVote.Encode()
