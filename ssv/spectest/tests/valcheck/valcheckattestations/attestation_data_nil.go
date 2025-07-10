package valcheckattestations

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BeaconVoteDataNil tests consensus data != nil
func BeaconVoteDataNil() tests.SpecTest {
	consensusData := &types.BeaconVote{
		Source: nil,
		Target: nil,
	}
	input, _ := consensusData.Encode()

	return &valcheck.SpecTest{
		Name:                "consensus data value check nil",
		Network:             types.PraterNetwork,
		RunnerRole:          types.RoleCommittee,
		DutySlot:            testingutils.TestingDutySlot,
		Input:               input,
		ExpectedSourceEpoch: 0,
		ExpectedTargetEpoch: 1,
		ExpectedError:       "attestation data source >= target",
	}
}
