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

	accdAggDataBytesF := func(cd *types.AggregatorCommitteeConsensusData) []byte {
		cdCopy := types.AggregatorCommitteeConsensusData{}
		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, &cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Aggregators[0].ValidatorIndex = testingutils.TestingWrongValidatorIndex

		ret, _ := cdCopy.Encode()
		return ret
	}
	accdSCCDataBytesF := func(cd *types.AggregatorCommitteeConsensusData) []byte {
		cdCopy := types.AggregatorCommitteeConsensusData{}
		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, &cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Contributors[0].ValidatorIndex = testingutils.TestingWrongValidatorIndex

		ret, _ := cdCopy.Encode()
		return ret
	}
	accdMixedDataBytesF := func(cd *types.AggregatorCommitteeConsensusData) []byte {
		cdCopy := types.AggregatorCommitteeConsensusData{}
		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, &cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Aggregators[0].ValidatorIndex = testingutils.TestingWrongValidatorIndex
		cdCopy.Contributors[0].ValidatorIndex = testingutils.TestingWrongValidatorIndex

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErrCode := types.WrongValidatorIndexErrorCode
	return valcheck.NewMultiSpecTest(
		"wrong validator index",
		testdoc.ValCheckDutyWrongValidatorIndexDoc,
		[]*valcheck.SpecTest{
			{
				Name:           "committee",
				Network:        types.BeaconTestNetwork,
				RunnerRole:     types.RoleCommittee,
				Input:          testingutils.TestBeaconVoteByts,
				ExpectedSource: *testingutils.TestBeaconVote.Source,
				ExpectedTarget: *testingutils.TestBeaconVote.Target,
				// No error since input doesn't contain validator index
			},
			{
				Name:       "aggregator committee scc",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      accdSCCDataBytesF(testingutils.TestSyncCommitteeContributionConsensusData),
				// No error since input doesn't contain validator index
			},
			{
				Name:       "aggregator committee agg",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      accdAggDataBytesF(testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)),
				// No error since input doesn't contain validator index
			},
			{
				Name:       "aggregator committee mixed",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      accdMixedDataBytesF(testingutils.TestAggregatorCommitteeConsensusDataForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(spec.DataVersionElectra), spec.DataVersionElectra)),
				// No error since input doesn't contain validator index
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
