package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "5f8dfc0b340c2fbac1c45d017eed5daeacacaba724a83006bf75029fc23f17c5",
	}
}
