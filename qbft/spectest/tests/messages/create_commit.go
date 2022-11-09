package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreateCommit tests creating a commit msg
func CreateCommit() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateCommit,
		Name:         "create commit",
		Value:        &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		Round:        10,
		ExpectedRoot: "4a58b7937892cfb0821c34e9fac161c982f3358c0dd4ff6b0d11cb9a455913cd",
	}
}
