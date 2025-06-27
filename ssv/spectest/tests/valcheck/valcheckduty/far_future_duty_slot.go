package valcheckduty

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FarFutureDutySlot tests duty.Slot higher than expected
func FarFutureDutySlot() tests.SpecTest {
	consensusDataBytsF := func(cd *types.ValidatorConsensusData) []byte {
		cdCopy := &types.ValidatorConsensusData{}

		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.Slot = 100000000

		ret, _ := cdCopy.Encode()
		return ret
	}

	expectedErr := "duty invalid: duty epoch is into far future"
	return valcheck.NewMultiSpecTest(
		"far future duty slot",
		"Tests duty value check with duty slot too far in the future",
		[]*valcheck.SpecTest{
			{
				Name:       "committee",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Input:      testingutils.TestBeaconVoteByts,
				// No error since input doesn't contain slot
			},
			{
				Name:          "sync committee aggregator",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestSyncCommitteeContributionConsensusData),
				ExpectedError: expectedErr,
			},
			{
				Name:          "aggregator phase0",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)),
				ExpectedError: expectedErr,
			},
			{
				Name:          "aggregator electra",
				Network:       types.BeaconTestNetwork,
				RunnerRole:    types.RoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData(spec.DataVersionElectra)),
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
	)
}
