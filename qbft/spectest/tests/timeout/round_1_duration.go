package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// Round1Duration tests timeout duration for round 1
func Round1Duration() *tests.MultiSpecTest {
	var testingNetwork = types.HoleskyNetwork
	var currentTime = testingNetwork.MinGenesisTime()

	return &tests.MultiSpecTest{
		Name: "round 1",
		Tests: []tests.SpecTest{
			&TimeoutDurationTest{
				Name:             "sync committee",
				Role:             types.BNRoleSyncCommittee,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      currentTime,
				ExpectedDuration: 6,
			},
			&TimeoutDurationTest{
				Name:             "sync committee contribution",
				Role:             types.BNRoleSyncCommitteeContribution,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      currentTime,
				ExpectedDuration: 10,
			},
			&TimeoutDurationTest{
				Name:             "attester",
				Role:             types.BNRoleAttester,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      currentTime,
				ExpectedDuration: 6,
			},
			&TimeoutDurationTest{
				Name:             "aggregator",
				Role:             types.BNRoleAggregator,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      currentTime,
				ExpectedDuration: 10,
			},
			&TimeoutDurationTest{
				Name:             "block proposer",
				Role:             types.BNRoleProposer,
				Height:           0,
				Round:            qbft.FirstRound,
				Network:          testingNetwork,
				CurrentTime:      currentTime,
				ExpectedDuration: 2,
			},
		},
	}

}
