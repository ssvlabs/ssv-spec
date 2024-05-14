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
	consensusDataBytsF := func(cd *types.ConsensusData) []byte {
		cdCopy := &types.ConsensusData{}

		b, _ := json.Marshal(cd)
		if err := json.Unmarshal(b, cdCopy); err != nil {
			panic(err.Error())
		}
		cdCopy.Duty.Slot = 100000000

		ret, _ := cdCopy.Encode()
		return ret
	}

	return &valcheck.MultiSpecTest{
		Name: "far future duty slot",
		Tests: []*valcheck.SpecTest{
			{
				Name:          "sync committee aggregator",
				Network:       types.BeaconTestNetwork,
				Role:          types.RoleSyncCommitteeContribution,
				Input:         consensusDataBytsF(testingutils.TestSyncCommitteeContributionConsensusData),
				ExpectedError: "duty invalid: duty epoch is into far future",
			},
			{
				Name:          "aggregator",
				Network:       types.BeaconTestNetwork,
				Role:          types.RoleAggregator,
				Input:         consensusDataBytsF(testingutils.TestAggregatorConsensusData),
				ExpectedError: "duty invalid: duty epoch is into far future",
			},
			{
				Name:          "proposer",
				Network:       types.BeaconTestNetwork,
				Role:          types.RoleProposer,
				Input:         consensusDataBytsF(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb)),
				ExpectedError: "duty invalid: duty epoch is into far future",
			},
			{
				Name:    "committee",
				Network: types.BeaconTestNetwork,
				Role:    types.RoleCommittee,
				Input:   testingutils.TestBeaconVoteByts,
				// No error since beacon vote doesn't include slot
			},
		},
	}
}
