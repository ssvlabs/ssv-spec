package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "fe85e25b3cf7168e9e6417a9daaa71567a3c0c689b633d8154e252d8225c113c",
	}
}
