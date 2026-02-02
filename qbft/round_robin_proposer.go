package qbft

import (
	"sort"

	"github.com/ssvlabs/ssv-spec/types"
)

// RoundRobinProposer returns the proposer for the round.
// Each new height starts with the first proposer and increments by 1 with each following round.
// Each new height has a different first round proposer which is +1 from the previous height.
// Also, the current Ethereum epoch is taken into account to introduce variability through epochs
// (mostly for committees with 4 operators, as 32%4 = 0 as the epochs would "repeat" otherwise).
// First height starts with index 0
func RoundRobinProposer(state *State, round Round) types.OperatorID {
	committee := sortedCommittee(state.CommitteeMember.Committee)

	firstRoundIndex := 0
	if state.Height != FirstHeight {
		firstRoundIndex += int(state.Height) % len(committee)
	}
	ethEpoch := int(state.Height) / 32

	index := (firstRoundIndex + int(round) - int(FirstRound) + ethEpoch) % len(committee)
	return committee[index].OperatorID
}

// sortedCommittee returns the original committee if it's already sorted by OperatorID.
// Otherwise it returns a sorted copy (leaving the original untouched).
func sortedCommittee(committee []*types.Operator) []*types.Operator {
	if sort.SliceIsSorted(committee, func(i, j int) bool {
		return committee[i].OperatorID < committee[j].OperatorID
	}) {
		return committee
	}

	sortedList := make([]*types.Operator, len(committee))
	copy(sortedList, committee)

	sort.Slice(sortedList, func(i, j int) bool {
		return sortedList[i].OperatorID < sortedList[j].OperatorID
	})
	return sortedList
}
