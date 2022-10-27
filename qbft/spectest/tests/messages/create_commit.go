package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		ValueRoot:    [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "d83735f8a99fabde910e4dff198163811be43ffbb1ffc62a80c75c44618270c1",
	}
}
