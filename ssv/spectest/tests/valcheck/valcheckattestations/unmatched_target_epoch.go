package valcheckattestations

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnmatchedEpochs tests AttestationData.Target.Epoch unmatched with expected
func UnmatchedTargetEpoch() tests.SpecTest {
	data := types.BeaconVote{
		BlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		Source: &spec.Checkpoint{
			Epoch: 0,
			Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Target: &spec.Checkpoint{
			Epoch: 2,
			Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
	}

	input, _ := data.Encode()

	return &valcheck.SpecTest{
		Name:                "attestation value check unmatched target epoch",
		Network:             types.BeaconTestNetwork,
		RunnerRole:          types.RoleCommittee,
		DutySlot:            testingutils.TestingDutySlot,
		Input:               input,
		ExpectedSourceEpoch: 1,
		ExpectedTargetEpoch: 2,
		ExpectedError:       "attestation data source/target epoch does not match expected",
	}
}
