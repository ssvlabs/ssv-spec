package testingutils

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// SortMapByKey sorts a map by its keys and returns a sorted slice of key-value pairs.
func SortedMapKeys[K constraints.Ordered, V any](m map[K]V) []struct {
	Key   K
	Value V
} {
	// Extract and sort keys
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// Create sorted key-value pairs
	pairs := make([]struct {
		Key   K
		Value V
	}, len(keys))

	for i, k := range keys {
		pairs[i] = struct {
			Key   K
			Value V
		}{k, m[k]}
	}

	return pairs
}
