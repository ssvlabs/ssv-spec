package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		ExpectedRoot: "0000000000000000000000000000000000000000000000000000000000000000",
	}
}
