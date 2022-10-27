package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		ValueRoot:    [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "aca8a000083a7f5a5cf0f7d011e857700ece644f095c5dae875dcda6ade672d1",
	}
}
