package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "284685effb1ce6e164cf990dc119119e9a7c8a765ca974faa7a820b23610e595",
	}
}
