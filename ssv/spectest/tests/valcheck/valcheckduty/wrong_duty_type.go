package valcheckduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDutyType tests duty.Type not attester
func WrongDutyType() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ValidatorConsensusData) []byte {
		input, _ := cd.Encode()
		return input
	}

	return &valcheck.MultiSpecTest{
		Name: "wrong duty type",
		Tests: []*valcheck.SpecTest{
			{
				Name:       "committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Input:      testingutils.TestBeaconVoteByts,
				// No error since input doesn't contain duty type
			},
			{
				Name:          "sync committee aggregator",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "proposer",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleProposer,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
		},
	}
}
