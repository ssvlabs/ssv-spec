package valcheckcommittee

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Slashable tests a slashable BeaconVote
func Slashable() tests.SpecTest {
	bv := &types.BeaconVote{
		BlockRoot: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		Source: &phase0.Checkpoint{
			Epoch: 0,
			Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Target: &phase0.Checkpoint{
			Epoch: 1,
			Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
	}
	input, _ := bv.Encode()

	attestationData := &phase0.AttestationData{
		Slot: testingutils.TestingDutySlot,
		// CommitteeIndex doesn't matter for slashing checks
		Index:           0,
		BeaconBlockRoot: bv.BlockRoot,
		Source:          nil,
		Target:          nil,
	}

	slashableRoot, _ := attestationData.HashTreeRoot()

	return &valcheck.SpecTest{
		Name:          "committee value check slashable",
		Network:       types.BeaconTestNetwork,
		Role:          types.RoleCommittee,
		Input:         input,
		ExpectedError: "slashable attestation",
		SlashableDataRoots: [][]byte{
			slashableRoot[:],
		},
	}
}
