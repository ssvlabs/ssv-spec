package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// Round1DurationFutureCurrentTime tests timeout duration for round 1 where the current time is after the expected
// start of the round
func Round1DurationFutureCurrentTime() tests.SpecTest {
	var testingNetwork = types.HoleskyNetwork
	height := qbft.FirstHeight
	round := qbft.FirstRound
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &MultiSpecTest{
		Name: "round 1 duration future current time",
		Tests: []*TimeoutDurationTest{
			{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 0,
			},
			{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 0,
			},
			{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 0,
			},
			{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12,
				ExpectedDuration: 0,
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
