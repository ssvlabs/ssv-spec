package valcheckduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
				BeaconRole:    types.BNRoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix)),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
			{
				Name:          "sync committee",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleSyncCommittee,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix)),
				ExpectedError: "duty invalid: wrong beacon role type", // it passes ConsensusData validation since  SyncCommitteeBlockRoot can't be nil, it's [32]byte
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix)),
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
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix)),
				ExpectedError: "duty invalid: wrong beacon role type",
			},
		},
	}
}
