package valcheckduty

import (
	"encoding/json"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongValidatorPK tests duty.PubKey wrong
func WrongValidatorPK() *valcheck.MultiSpecTest {
	consensusDataBytsF := func(cd *types.ConsensusData) []byte {
		cdCopy := &types.ConsensusData{}
		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.PubKey = testingutils.TestingWrongValidatorPubKey

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErr := "duty invalid: wrong validator pk"
	return &valcheck.MultiSpecTest{
		Name: "wrong validator PK",
		Tests: []*valcheck.SpecTest{
			{
				Name:          "sync committee aggregator",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestSyncCommitteeContributionConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "sync committee",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleSyncCommittee,
				Input:         consensusDataBytsF(testingutils.TestSyncCommitteeConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "proposer",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleProposer,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "attester",
				Network:       types.BeaconTestNetwork,
				BeaconRole:    types.BNRoleAttester,
				Input:         consensusDataBytsF(testingutils.TestAttesterConsensusData),
				ExpectedError: expectedErr,
			},
		},
	}

}
