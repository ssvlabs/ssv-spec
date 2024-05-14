package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "d6a58346ff2236d3c5e818d8b9d825a879e78886d895e5b645a20a22d7f50cbb",
	}
}
