package valcheckduty

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongValidatorIndex tests duty.ValidatorIndex wrong
func WrongValidatorIndex() *valcheck.MultiSpecTest {
	consensusDataBytsF := func(role types.BeaconRole) *qbft.Data {
		cd := &types.ConsensusData{
			Duty: &types.Duty{
				Type:                    role,
				PubKey:                  testingutils.TestingValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex + 100,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
			AttestationData: &phase0.AttestationData{
				Slot:            1000,
				Index:           3,
				BeaconBlockRoot: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
				Source: &phase0.Checkpoint{
					Epoch: 0,
					Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
				},
			},
			BlockData:              testingutils.TestingBeaconBlock,
			AggregateAndProof:      testingutils.TestingAggregateAndProof,
			SyncCommitteeBlockRoot: testingutils.TestingSyncCommitteeBlockRoot,
			SyncCommitteeContribution: types.ContributionsMap{
				testingutils.TestingContributionProofsSigned[0]: testingutils.TestingSyncCommitteeContributions[0],
				testingutils.TestingContributionProofsSigned[1]: testingutils.TestingSyncCommitteeContributions[1],
				testingutils.TestingContributionProofsSigned[2]: testingutils.TestingSyncCommitteeContributions[2],
			},
		}

		source, _ := cd.MarshalSSZ()
		root, _ := cd.HashTreeRoot()
		return &qbft.Data{
			Root:   root,
			Source: source,
		}
	}

	return &valcheck.MultiSpecTest{
		Name: "wrong validator index",
		Tests: []*valcheck.SpecTest{
			{
				Name:          "sync committee aggregator",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(types.BNRoleSyncCommitteeContribution),
				ExpectedError: "duty invalid: wrong validator index",
			},
			{
				Name:          "sync committee",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleSyncCommittee,
				Input:         consensusDataBytsF(types.BNRoleSyncCommittee),
				ExpectedError: "duty invalid: wrong validator index",
			},
			{
				Name:          "aggregator",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleAggregator,
				Input:         consensusDataBytsF(types.BNRoleAggregator),
				ExpectedError: "duty invalid: wrong validator index",
			},
			{
				Name:          "proposer",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleProposer,
				Input:         consensusDataBytsF(types.BNRoleProposer),
				ExpectedError: "duty invalid: wrong validator index",
			},
			{
				Name:          "attester",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleAttester,
				Input:         consensusDataBytsF(types.BNRoleAttester),
				ExpectedError: "duty invalid: wrong validator index",
			},
		},
	}
}
