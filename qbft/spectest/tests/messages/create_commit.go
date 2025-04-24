package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "b483d818d745f22091a7360afe1e214bb3bf645375062de2c636fff1bace7c94",
	}
}
