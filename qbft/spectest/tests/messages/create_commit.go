package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "ce830104bddf6905db857b74df322663359c3c4bb31d8f0554b3ff4a78f96dda",
	}
}
