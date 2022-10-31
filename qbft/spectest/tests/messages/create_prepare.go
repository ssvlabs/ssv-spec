package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreatePrepare tests creating a prepare msg
func CreatePrepare() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		Round:        10,
		ExpectedRoot: "0102030400000000000000000000000000000000000000000000000000000000",
	}
}
