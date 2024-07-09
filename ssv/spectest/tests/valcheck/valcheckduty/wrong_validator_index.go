package valcheckduty

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongValidatorIndex tests duty.ValidatorIndex wrong
func WrongValidatorIndex() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ValidatorConsensusData) []byte {
		cdCopy := &types.ValidatorConsensusData{}
		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.ValidatorIndex = testingutils.TestingValidatorIndex + 100

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErr := "duty invalid: wrong validator index"
	return &valcheck.MultiSpecTest{
		Name: "wrong validator index",
		Tests: []*valcheck.SpecTest{
			{
				Name:       "committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Input:      testingutils.TestBeaconVoteByts,
				// No error since input doesn't contain validator index
			},
			{
				Name:          "sync committee aggregator",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestSyncCommitteeContributionConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "proposer",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleProposer,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedError: expectedErr,
			},
		},
	}
}
