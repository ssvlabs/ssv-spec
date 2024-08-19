package testingutils

import "math/rand"

// Function to checks if all given lists are empty
func allListsEmpty(lists [][]interface{}) bool {
	for _, list := range lists {
		if len(list) > 0 {
			return false
		}
	}
	return true
}

// Function to return the indices of non-empty lists
func getNonEmptyListsIndices(lists [][]interface{}) []int {
	indices := make([]int, 0)
	for i, list := range lists {
		if len(list) > 0 {
			indices = append(indices, i)
		}
	}
	return indices
}

// Function to merge the given lists with a random pick order.
// The result is a shuffled list that preserves the orders of each sublist
func MergeListsWithRandomPick(lists [][]interface{}) []interface{} {

	seed := int64(42)
	rand := rand.New(rand.NewSource(seed))
	result := make([]interface{}, 0)

	// Continue until all lists are empty
	for !allListsEmpty(lists) {
		// Select a random non-empty list
		nonEmptyListsIndices := getNonEmptyListsIndices(lists)
		if len(nonEmptyListsIndices) == 0 {
			break
		}
		chosenListIndex := nonEmptyListsIndices[rand.Intn(len(nonEmptyListsIndices))]

		// Append the next element from the chosen list to the result
		result = append(result, lists[chosenListIndex][0])

		// Remove the element from the chosen list
		lists[chosenListIndex] = lists[chosenListIndex][1:]
	}

	return result
}
