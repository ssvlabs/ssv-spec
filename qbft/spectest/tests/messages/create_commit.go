package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "6650657e621029ea09937d8f2766b3ea8b71ec19fc02201d9c9c9a45c2c1baa8",
	}
}
