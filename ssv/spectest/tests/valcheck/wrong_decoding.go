package valcheck

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDecoding tests on bad decoding case
func WrongDecoding() tests.SpecTest {

	return NewMultiSpecTest(
		"wrong decoding",
		testdoc.ValCheckBadDecodingDoc,
		[]*SpecTest{
			{
				Name:              "committee",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleCommittee,
				Input:             testingutils.TestProposerConsensusDataBytsV(spec.DataVersionDeneb),
				ExpectedSource:    *testingutils.TestBeaconVote.Source,
				ExpectedTarget:    *testingutils.TestBeaconVote.Target,
				ExpectedErrorCode: types.DecodeBeaconVoteErrorCode,
			},
			{
				Name:              "aggregator committee phase0",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleAggregatorCommittee,
				Input:             testingutils.TestBeaconVoteByts,
				ExpectedErrorCode: types.AggCommConsensusDataDecodeErrorCode,
			},
			{
				Name:              "aggregator committee electra",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleAggregatorCommittee,
				Input:             testingutils.TestBeaconVoteByts,
				ExpectedErrorCode: types.AggCommConsensusDataDecodeErrorCode,
			},
			{
				Name:              "proposer",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             testingutils.TestBeaconVoteByts,
				ExpectedErrorCode: types.ProposerConsensusDataDecodeErrorCode,
			},
		},
	)
}
