package valcheckduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongValidatorPK tests duty.PubKey wrong
func WrongValidatorPK() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ProposerConsensusData) []byte {
		b, err := cd.Encode()
		if err != nil {
			panic(err.Error())
		}
		cdCopy := &types.ProposerConsensusData{}
		if err := cdCopy.Decode(b); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.PubKey = testingutils.TestingWrongValidatorPubKey

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErrCode := types.WrongValidatorPubkeyErrorCode
	return valcheck.NewMultiSpecTest(
		"wrong validator PK",
		testdoc.ValCheckDutyWrongValidatorPKDoc,
		[]*valcheck.SpecTest{
			{
				Name:           "committee",
				Network:        types.BeaconTestNetwork,
				RunnerRole:     types.RoleCommittee,
				Input:          testingutils.TestBeaconVoteByts,
				ExpectedSource: *testingutils.TestBeaconVote.Source,
				ExpectedTarget: *testingutils.TestBeaconVote.Target,
				// No error since input doesn't contain validator public key
			},
			{
				Name:       "aggregator committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(gloas.DataVersionGloas), gloas.DataVersionGloas),
				// No error since input doesn't contain validator public key
			},
			{
				Name:              "proposer",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             consensusDataBytsF(testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas)),
				ExpectedErrorCode: expectedErrCode,
			},
		},
	)
}
