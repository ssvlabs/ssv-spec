package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "22763e542ac1faba4c01c9aefee5b96dc23889c1fb079e46e00a13c9fcc367c7",
	}
}
