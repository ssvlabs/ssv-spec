package valcheckduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FarFutureDutySlot tests duty.Slot higher than expected
func FarFutureDutySlot() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ProposerConsensusData) []byte {
		// Deep-copy via SSZ, not JSON: the Gloas placeholder version is outside go-eth2-client's
		// enum, and spec.DataVersion.MarshalJSON panics on out-of-enum values.
		b, err := cd.Encode()
		if err != nil {
			panic(err.Error())
		}
		cdCopy := &types.ProposerConsensusData{}
		if err := cdCopy.Decode(b); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.Slot = 100000000

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErrCode := types.DutyEpochTooFarFutureErrorCode
	return valcheck.NewMultiSpecTest(
		"far future duty slot",
		testdoc.ValCheckDutyFarFutureDutySlotDoc,
		[]*valcheck.SpecTest{
			{
				Name:           "committee",
				Network:        types.BeaconTestNetwork,
				RunnerRole:     types.RoleCommittee,
				Input:          testingutils.TestBeaconVoteByts,
				ExpectedSource: *testingutils.TestBeaconVote.Source,
				ExpectedTarget: *testingutils.TestBeaconVote.Target,
				// No error since input doesn't contain slot
			},
			{
				Name:       "aggregator committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(gloas.DataVersionGloas), gloas.DataVersionGloas),
				// No error since input doesn't contain slot
			},
			{
				Name:       "aggregator committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleAggregatorCommittee,
				Input:      testingutils.TestAggregatorCommitteeConsensusDataBytesForDuty(testingutils.TestingAggregatorCommitteeDutyMixed(gloas.DataVersionGloas), gloas.DataVersionGloas),
				// No error since input doesn't contain slot
			},
			{
				// The value must be Gloas-versioned — another Version would trip the version guard before
				// the far-future check this test pins (that guard has its own vector in valcheckproposer).
				Name:              "proposer",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             consensusDataBytsF(testingutils.TestProposerConsensusDataV(gloas.DataVersionGloas)),
				ExpectedErrorCode: expectedErrCode,
			},
		},
	)
}
