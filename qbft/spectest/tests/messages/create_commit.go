package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:      tests.CreateCommit,
		Name:            "create commit",
		Value:           []byte{1, 2, 3, 4},
		Round:           10,
		ExpectedSSZRoot: "a30c2e98de75d97701f3d32026b6a6b67df7bac2aa2804ef7ea481ca03542850",
	}
}
