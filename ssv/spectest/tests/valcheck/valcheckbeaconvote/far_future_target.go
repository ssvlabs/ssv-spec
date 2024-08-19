package valcheckbeaconvote

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FarFutureTarget tests AttestationData.Target.Epoch higher than expected
func FarFutureTarget() tests.SpecTest {
	data := types.BeaconVote{
		BlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		Source: &spec.Checkpoint{
			Epoch: 0,
			Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Target: &spec.Checkpoint{
			Epoch: 10000000,
			Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
	}

	input, _ := data.Encode()

	return &valcheck.MultiSpecTest{
		Name: "beacon vote value check far future target",
		Tests: []*valcheck.SpecTest{
			{
				Name:             "attestation duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingAttesterDuty,
				Input:            input,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
				ExpectedError:    "attestation data target epoch is into far future",
			},
			{
				Name:             "sync committee duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingSyncCommitteeDuty,
				Input:            input,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
				ExpectedError:    "attestation data target epoch is into far future",
			},
			{
				Name:             "attestation and sync committee duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingAttesterAndSyncCommitteeDuties,
				Input:            input,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
				ExpectedError:    "attestation data target epoch is into far future",
			},
		},
	}
}
