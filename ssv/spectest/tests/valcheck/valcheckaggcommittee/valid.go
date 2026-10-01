package valcheckaggcommittee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests the valid scenario
func Valid() tests.SpecTest {

	return valcheck.NewMultiSpecTest(
		"aggcommittee value check valid",
		testdoc.ValCheckAggCommitteeValidDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "aggregator",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyOnlyAggregator(gloas.DataVersionGloas), gloas.DataVersionGloas),
			},
			{
				Name:       "sync committee contribution",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyOnlySyncCommittee(), gloas.DataVersionGloas),
			},
			{
				Name:       "mixed",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(gloas.DataVersionGloas), gloas.DataVersionGloas),
			},
		},
	)
}
