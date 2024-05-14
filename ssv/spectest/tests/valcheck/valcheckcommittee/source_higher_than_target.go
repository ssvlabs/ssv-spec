package valcheckcommittee

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
)

// SourceHigherThanTarget tests a beacon vote with the source higher than the target
func SourceHigherThanTarget() tests.SpecTest {
	beaconVote := &types.BeaconVote{
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

	input, _ := beaconVote.Encode()

	return &valcheck.SpecTest{
		Name:          "committee value check source higher than target",
		Network:       types.BeaconTestNetwork,
		Role:          types.RoleCommittee,
		Input:         input,
		ExpectedError: "attestation data source > target",
	}
}
