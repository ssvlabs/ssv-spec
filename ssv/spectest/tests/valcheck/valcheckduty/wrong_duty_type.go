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
	consensusDataBytsF := func(cd *types.ConsensusData) []byte {
		input, _ := cd.Encode()
		return input
	}

	return &valcheck.MultiSpecTest{
		Name: "wrong duty type",
		Tests: []*valcheck.SpecTest{
			{
				Name:          "sync committee aggregator",
				Network:       types.BeaconTestNetwork,
				Role:          types.RoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				Role:          types.RoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "proposer",
				Network:       types.BeaconTestNetwork,
				Role:          types.RoleProposer,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:    "committee",
				Network: types.BeaconTestNetwork,
				Role:    types.RoleCommittee,
				Input:   testingutils.TestBeaconVoteByts,
				// No error since beacon vote doesn't include role
			},
		},
	}
}
