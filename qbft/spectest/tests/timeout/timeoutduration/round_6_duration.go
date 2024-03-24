package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// Round6Duration tests timeout duration for round 8 where the current time is the expected start of the round
func Round6Duration() tests.SpecTest {
	testingNetwork := types.HoleskyNetwork
	height := qbft.FirstHeight
	var round qbft.Round = 6
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &MultiSpecTest{
		Name: "round 6 duration",
		Tests: []*TimeoutDurationTest{
			{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 14,
				ExpectedDuration: 2,
			},
			{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 18,
				ExpectedDuration: 2,
			},
			{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 14,
				ExpectedDuration: 2,
			},
			{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 18,
				ExpectedDuration: 2,
			},
			{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 2,
			},
		},
	}

}
