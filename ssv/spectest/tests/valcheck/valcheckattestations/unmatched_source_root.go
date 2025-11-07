package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnmatchedSourceRoot tests AttestationData.Source.Root unmatched with expected
func UnmatchedSourceRoot() tests.SpecTest {
	data := types.BeaconVote{
		BlockRoot: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		Source: &phase0.Checkpoint{
			Epoch: 0,
			Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Target: &phase0.Checkpoint{
			Epoch: 2,
			Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
	}

	expectedSourceCheckpoint := phase0.Checkpoint{
		Epoch: 0,
		// different from above
		Root: phase0.Root{9, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	}
	expectedTargetCheckpoint := phase0.Checkpoint{
		Epoch: 2,
		Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
	}
	input, _ := data.Encode()

	return valcheck.NewSpecTest(
		"attestation value check unmatched source root",
		testdoc.ValCheckAttestationUnmatchedSourceRootDoc, // Documentation string if any is needed
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlot,
		input,
		expectedSourceCheckpoint,
		expectedTargetCheckpoint,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.CheckpointMismatch,
		false,
	)
}
