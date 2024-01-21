package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// Round5DurationOnRound3Time tests timeout duration for round 5 where the current time is the expected start of round 3
func Round5DurationOnRound3Time() tests.SpecTest {
	testingNetwork := types.HoleskyNetwork
	height := qbft.Height(10)
	var round qbft.Round = 5
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &tests.MultiSpecTest{
		Name: "round 5 duration on round 3 time",
		Tests: []tests.SpecTest{
			&TimeoutDurationTest{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 8,
				ExpectedDuration: 6,
			},
			&TimeoutDurationTest{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 6,
			},
			&TimeoutDurationTest{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 8,
				ExpectedDuration: 6,
			},
			&TimeoutDurationTest{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 6,
			},
			&TimeoutDurationTest{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 6,
				ExpectedDuration: 2,
			},
		},
	}

}
