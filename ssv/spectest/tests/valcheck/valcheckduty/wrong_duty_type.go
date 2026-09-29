package valcheckduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDutyType tests duty.Type not attester
func WrongDutyType() tests.SpecTest {
	consensusDataBytsWithIncorrectTypeF := func(cd *types.ProposerConsensusData) []byte {
		cd.Duty.Type = types.BNRoleAggregator
		input, _ := cd.Encode()
		return input
	}

	return valcheck.NewMultiSpecTest(
		"wrong duty type",
		testdoc.ValCheckDutyWrongDutyTypeDoc,
		[]*valcheck.SpecTest{
			{
				Name:           "committee",
				Network:        types.BeaconTestNetwork,
				RunnerRole:     types.RoleCommittee,
				Input:          testingutils.TestBeaconVoteByts,
				ExpectedSource: *testingutils.TestBeaconVote.Source,
				ExpectedTarget: *testingutils.TestBeaconVote.Target,
				// No error since input doesn't contain duty type
			},
			{
				Name:       "aggregator committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(gloas.DataVersionGloas), gloas.DataVersionGloas),
				// No error since input doesn't contain duty type
			},
			{
				Name:              "proposer",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             consensusDataBytsWithIncorrectTypeF(testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas)),
				ExpectedErrorCode: types.QBFTValueInvalidErrorCode,
			},
		},
	)
}
