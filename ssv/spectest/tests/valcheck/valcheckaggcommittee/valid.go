package valcheckaggcommittee

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests the valid scenario
func Valid() tests.SpecTest {

	return valcheck.NewMultiSpecTest(
		"aggcommittee value check valid",
		testdoc.ValCheckAggCommitteeValidDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "aggregator phase0",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyOnlyAggregator(spec.DataVersionPhase0), spec.DataVersionPhase0),
			},
			{
				Name:       "aggregator electra",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyOnlyAggregator(spec.DataVersionElectra), spec.DataVersionElectra),
			},
			{
				Name:       "sync committee contribution",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyOnlySyncCommittee(), spec.DataVersionElectra),
			},
			{
				Name:       "mixed",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(spec.DataVersionElectra), spec.DataVersionElectra),
			},
		},
	)
}
