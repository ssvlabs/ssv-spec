package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// Round13Duration tests timeout duration for round 13 where the current time is the expected start of the round
func Round13Duration() tests.SpecTest {
	testingNetwork := types.HoleskyNetwork
	height := qbft.Height(40)
	var round qbft.Round = 13
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &MultiSpecTest{
		Name: "round 13 duration",
		Tests: []*TimeoutDurationTest{
			{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 16 + 6*120,
				ExpectedDuration: 120,
			},
			{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 20 + 6*120,
				ExpectedDuration: 120,
			},
			{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 16 + 6*120,
				ExpectedDuration: 120,
			},
			{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 20 + 6*120,
				ExpectedDuration: 120,
			},
			{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           height,
				Round:            round,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime + 12 + 6*120,
				ExpectedDuration: 120,
			},
		},
	}

}
