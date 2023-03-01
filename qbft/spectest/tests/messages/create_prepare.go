package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "2eba5b18818e0ec94f5c02ff7abc8ca932ed5d1f32a115197fbaa14247a39cb2",
	}
}
