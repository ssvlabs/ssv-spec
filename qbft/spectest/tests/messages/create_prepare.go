package messages

import "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "b6404cc0c4dbdc2fa8bfcbebcd75344dfa976da4a8c2487191aa34c8d64af441",
	}
}
