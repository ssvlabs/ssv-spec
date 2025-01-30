package qbft

import "github.com/ssvlabs/ssv-spec/types"

// RoundRobinProposer returns the proposer for the round.
// Each new height starts with the first proposer and increments by 1 with each following round.
// Each new height has a different first round proposer which is +1 from the previous height.
// First height starts with index 0
// Assumption: the operators list is sorted.
func RoundRobinProposer(height Height, round Round, operators []types.OperatorID) types.OperatorID {
	firstRoundIndex := 0
	if height != FirstHeight {
		firstRoundIndex += int(height) % len(operators)
	}

	index := (firstRoundIndex + int(round) - int(FirstRound)) % len(operators)
	return operators[index]
}
