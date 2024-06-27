package valcheckduty

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongValidatorPK tests duty.PubKey wrong
func WrongValidatorPK() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ValidatorConsensusData) []byte {
		cdCopy := &types.ValidatorConsensusData{}
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
				Name:       "committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Input:      testingutils.TestBeaconVoteByts,
				// No error since input doesn't contain validator public key
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
