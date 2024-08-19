package valcheckbeaconvote

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BeaconVoteDataNil tests consensus data != nil
func BeaconVoteDataNil() tests.SpecTest {
	consensusData := &types.BeaconVote{
		Source: nil,
		Target: nil,
	}
	input, _ := consensusData.Encode()

	return &valcheck.MultiSpecTest{
		Name: "beacon vote data value check nil",
		Tests: []*valcheck.SpecTest{
			{
				Name:             "attestation duty",
				Network:          types.PraterNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingAttesterDuty,
				Input:            input,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
				// TODO: due to decoding we get the wierd error.. if one of them are nils we may pass..
				// not sure how important to fix this...
				ExpectedError: "attestation data source >= target",
			},
			{
				Name:             "sync committee duty",
				Network:          types.PraterNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingSyncCommitteeDuty,
				Input:            input,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
				ExpectedError:    "attestation data source >= target",
			},
			{
				Name:             "attestation and sync committee duty",
				Network:          types.PraterNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingAttesterAndSyncCommitteeDuties,
				Input:            input,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
				ExpectedError:    "attestation data source >= target",
			},
		},
	}
}
