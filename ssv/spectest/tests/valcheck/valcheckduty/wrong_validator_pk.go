package valcheckduty

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongValidatorPK tests duty.PubKey wrong
func WrongValidatorPK() *valcheck.MultiSpecTest {
	consensusDataBytsF := func(role types.BeaconRole) []byte {
		data := &types.ConsensusData{
			Duty: &types.Duty{
				Type:                    role,
				PubKey:                  testingutils.TestingWrongValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
			AttestationData: &spec.AttestationData{
				Slot:            1000,
				Index:           3,
				BeaconBlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
				Source: &spec.Checkpoint{
					Epoch: 0,
					Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
				},
				Target: &spec.Checkpoint{
					Epoch: 1,
					Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
				},
			},
		}

		input, _ := data.Encode()
		return input
	}

	return &valcheck.MultiSpecTest{
		Name: "wrong validator PK",
		Tests: []*valcheck.SpecTest{
			{
				Name:          "sync committee aggregator",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(types.BNRoleSyncCommitteeContribution),
				ExpectedError: "duty invalid: wrong validator pk",
			},
			{
				Name:          "sync committee",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleSyncCommittee,
				Input:         consensusDataBytsF(types.BNRoleSyncCommittee),
				ExpectedError: "duty invalid: wrong validator pk",
			},
			{
				Name:          "aggregator",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleAggregator,
				Input:         consensusDataBytsF(types.BNRoleAggregator),
				ExpectedError: "duty invalid: wrong validator pk",
			},
			{
				Name:          "proposer",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleProposer,
				Input:         consensusDataBytsF(types.BNRoleProposer),
				ExpectedError: "duty invalid: wrong validator pk",
			},
			{
				Name:          "attester",
				Network:       types.NowTestNetwork,
				BeaconRole:    types.BNRoleAttester,
				Input:         consensusDataBytsF(types.BNRoleAttester),
				ExpectedError: "duty invalid: wrong validator pk",
			},
		},
	}

}
