package valcheckduty

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongDutyType tests duty.Type not attester
func WrongDutyType() *valcheck.MultiSpecTest {
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
				BeaconRole:    types.BNRoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "sync committee",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleSyncCommittee,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type", // it passes ConsensusData validation since  SyncCommitteeBlockRoot can't be nil, it's [32]byte
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "proposer",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleProposer,
				Input:         consensusDataBytsF(testingutils.TestAttesterConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "attester",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleAttester,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusData),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
		},
	}
}
