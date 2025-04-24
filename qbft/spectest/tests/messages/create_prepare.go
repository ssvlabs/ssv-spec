package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "b7469add09080fcaaa9e702509c39b1bbc0d0485f8272a3e20d400ce13fc6d60",
	}
}
