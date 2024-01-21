package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/timeout"
	"github.com/bloxapp/ssv-spec/types"
)

// Round14DurationOnRound13Time tests timeout duration for round 14 where the current time is the expected start of
// round 13
func Round14DurationOnRound13Time() *tests.MultiSpecTest {
	testingNetwork := types.HoleskyNetwork
	height := qbft.Height(40)
	var round qbft.Round = 14
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &tests.MultiSpecTest{
		Name: "round 14 duration on round 13 time",
		Tests: []tests.SpecTest{
			&timeout.TimeoutDurationTest{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 20 + 4*120,
				ExpectedDuration: 240,
			},
			&timeout.TimeoutDurationTest{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 24 + 4*120,
				ExpectedDuration: 240,
			},
			&timeout.TimeoutDurationTest{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 20 + 4*120,
				ExpectedDuration: 240,
			},
			&timeout.TimeoutDurationTest{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 24 + 4*120,
				ExpectedDuration: 240,
			},
			&timeout.TimeoutDurationTest{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 20 + 4*120,
				ExpectedDuration: 120,
			},
		},
	}

}
