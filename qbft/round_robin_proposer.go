package qbft

import "github.com/ssvlabs/ssv-spec/types"

// RoundRobinProposer returns the proposer for the round.
// Each new height starts with the first proposer and increments by 1 with each following round.
// Each new height has a different first round proposer which is +1 from the previous height.
// First height starts with index 0
func RoundRobinProposer(state *State, round Round) types.OperatorID {
	committeeLength := uint64(len(state.CommitteeMember.Committee))

	firstRoundIndex := uint64(0)

	// Increment index using height
	if state.Height != FirstHeight {
		firstRoundIndex += uint64(state.Height) % committeeLength
	}

	// Increment index using round heights
	index := (firstRoundIndex + uint64(round) - uint64(FirstRound)) % committeeLength
	return state.CommitteeMember.Committee[index].OperatorID
}
