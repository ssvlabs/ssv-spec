package qbft

import "github.com/ssvlabs/ssv-spec/types"

// RoundRobinProposer returns the proposer for the round.
// Each new height starts with the first proposer and increments by 1 with each following round.
// Each new height has a different first round proposer which is +1 from the previous height.
// First height starts with index 0
func RoundRobinProposer(state *State, round Round) types.OperatorID {
	firstRoundIndex := 0
	if state.Height != FirstHeight {
		firstRoundIndex += int(state.Height) % len(state.CommitteeMember.Committee)
	}

	index := (firstRoundIndex + int(round) - int(FirstRound)) % len(state.CommitteeMember.Committee)
	return state.CommitteeMember.Committee[index].OperatorID
}
