package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests/timeout"
	"github.com/bloxapp/ssv-spec/types"
)

// Round8Duration tests timeout duration for round 8 where the current time is the expected start of the round
func Round8Duration() *tests.MultiSpecTest {
	testingNetwork := types.HoleskyNetwork
	height := qbft.FirstHeight
	var round qbft.Round = 8
	dutyStartTime := uint64(testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height)))

	return &tests.MultiSpecTest{
		Name: "round 8 duration",
		Tests: []tests.SpecTest{
			&timeout.TimeoutDurationTest{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 18,
				ExpectedDuration: 2,
			},
			&timeout.TimeoutDurationTest{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 22,
				ExpectedDuration: 2,
			},
			&timeout.TimeoutDurationTest{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 18,
				ExpectedDuration: 2,
			},
			&timeout.TimeoutDurationTest{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 22,
				ExpectedDuration: 2,
			},
			&timeout.TimeoutDurationTest{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 16,
				ExpectedDuration: 2,
			},
		},
	}

}
