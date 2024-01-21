package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// Round3DurationOnDutyStartTime tests timeout duration for round 3 where the current time is the start of the duty
func Round3DurationOnDutyStartTime() tests.SpecTest {
	testingNetwork := types.HoleskyNetwork
	height := qbft.FirstHeight
	var round qbft.Round = 3
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &tests.MultiSpecTest{
		Name: "round 3 on duty start time",
		Tests: []tests.SpecTest{
			&TimeoutDurationTest{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 10,
			},
			&TimeoutDurationTest{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 14,
			},
			&TimeoutDurationTest{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 10,
			},
			&TimeoutDurationTest{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 14,
			},
			&TimeoutDurationTest{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 2,
			},
		},
	}

}
