package valcheckcommittee

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
)

// FarFutureTarget tests a beacon vote with a target higher than expected
func FarFutureTarget() tests.SpecTest {
	beaconVote := &types.BeaconVote{
		BlockRoot: phase0.Root{1, 2, 3, 4},
		Source: &phase0.Checkpoint{
			Epoch: 0,
			Root:  phase0.Root{1, 2, 3, 4},
		},
		Target: &phase0.Checkpoint{
			Epoch: 10000000,
			Root:  phase0.Root{1, 2, 3, 5},
		},
	}

	input, _ := beaconVote.Encode()

	return &valcheck.SpecTest{
		Name:          "committee value check far future target",
		Network:       types.BeaconTestNetwork,
		Role:          types.RoleCommittee,
		Input:         input,
		ExpectedError: "attestation data target epoch is into far future",
	}
}
