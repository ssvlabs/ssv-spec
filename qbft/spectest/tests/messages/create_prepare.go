package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "0d1ed244a454b3df5e80fac60ed426ee1d3166534b35cd2f6b1602cddef42ded",
	}
}
