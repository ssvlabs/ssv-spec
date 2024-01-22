package timeoutduration

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"

	"github.com/bloxapp/ssv-spec/types"
)

// Round1Duration tests timeout duration for round 1
func Round1Duration() tests.SpecTest {
	var testingNetwork = types.HoleskyNetwork
	height := qbft.FirstHeight
	dutyStartTime := testingNetwork.EstimatedTimeAtSlot(phase0.Slot(height))

	return &MultiSpecTest{
		Name: "round 1",
		Tests: []*TimeoutDurationTest{
			{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 6,
			},
			{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 10,
			},
			{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 6,
			},
			{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 10,
			},
			{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      dutyStartTime,
				ExpectedDuration: 2,
			},
		},
	}

}
