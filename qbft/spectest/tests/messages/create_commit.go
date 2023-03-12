package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "0834b51f3c87d4aba362d7e2eeb4172d22e0ef18d4dfadd37e8c9ceb62c7719d",
	}
}
