package valcheckduty

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongValidatorIndex tests duty.ValidatorIndex wrong
func WrongValidatorIndex() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ValidatorConsensusData) []byte {
		cdCopy := types.ValidatorConsensusData{}
		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, &cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.ValidatorIndex = testingutils.TestingWrongValidatorIndex

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErrCode := types.WrongValidatorIndexErrorCode
	return valcheck.NewMultiSpecTest(
		"wrong validator index",
		testdoc.ValCheckDutyWrongValidatorIndexDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Input:      testingutils.TestBeaconVoteByts,
				// No error since input doesn't contain validator index
			},
			{
				Name:              "sync committee aggregator",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleSyncCommitteeContribution,
				Input:             consensusDataBytsF(testingutils.TestSyncCommitteeContributionConsensusData),
				ExpectedErrorCode: expectedErrCode,
			},
			{
				Name:              "aggregator phase0",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleAggregator,
				Input:             consensusDataBytsF(testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)),
				ExpectedErrorCode: expectedErrCode,
			},
			{
				Name:              "aggregator electra",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleAggregator,
				Input:             consensusDataBytsF(testingutils.TestAggregatorConsensusData(spec.DataVersionElectra)),
				ExpectedErrorCode: expectedErrCode,
			},
			{
				Name:              "proposer",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedErrorCode: expectedErrCode,
			},
		},
	)
}
