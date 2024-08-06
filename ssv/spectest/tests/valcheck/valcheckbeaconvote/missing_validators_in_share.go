package valcheckbeaconvote

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MissingValidatorInShare tests creation of a value check from a duty that has unknown validator shares
func MissingValidatorInShare() tests.SpecTest {
	return &valcheck.MultiSpecTest{
		Name: "beacon vote value check with missing validators in share",
		Tests: []*valcheck.SpecTest{
			{
				Name:       "attestation duty",
				Network:    types.PraterNetwork,
				RunnerRole: types.RoleCommittee,
				Duty: testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot,
					[]phase0.ValidatorIndex{5}),
				Input:            testingutils.TestBeaconVoteByts,
				ValidatorsShares: testingutils.Testing4ValidatorsShareMap,
				ExpectedError:    "assigned validator duty doesn't have a validator share",
			},
			{
				Name:       "sync committee duty",
				Network:    types.PraterNetwork,
				RunnerRole: types.RoleCommittee,
				Duty: testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot,
					[]phase0.ValidatorIndex{5}),
				Input:            testingutils.TestBeaconVoteByts,
				ValidatorsShares: testingutils.Testing4ValidatorsShareMap,
			},
			{
				Name:       "attestation and sync committee duty",
				Network:    types.PraterNetwork,
				RunnerRole: types.RoleCommittee,
				Duty: testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot,
					[]phase0.ValidatorIndex{5}, []phase0.ValidatorIndex{5}),
				Input:            testingutils.TestBeaconVoteByts,
				ValidatorsShares: testingutils.Testing4ValidatorsShareMap,
				ExpectedError:    "assigned validator duty doesn't have a validator share",
			},
		},
	}
}
