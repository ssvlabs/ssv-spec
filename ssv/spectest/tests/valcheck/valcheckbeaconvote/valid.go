package valcheckbeaconvote

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests valid data
func Valid() tests.SpecTest {
	return &valcheck.MultiSpecTest{
		Name: "beacon value check valid",
		Tests: []*valcheck.SpecTest{
			{
				Name:             "attestation duty",
				Network:          types.PraterNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingAttesterDuty,
				Input:            testingutils.TestBeaconVoteByts,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
			},
			{
				Name:             "sync committee duty",
				Network:          types.PraterNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingSyncCommitteeDuty,
				Input:            testingutils.TestBeaconVoteByts,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
			},
			{
				Name:             "attestation and sync committee duty",
				Network:          types.PraterNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(1), testingutils.ValidatorIndexList(1)),
				Input:            testingutils.TestBeaconVoteByts,
				ValidatorsShares: testingutils.TestingSingleValidatorShareMap,
			},
		},
	}
}
