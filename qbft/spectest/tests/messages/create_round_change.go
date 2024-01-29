package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "6a2199b09d7f31f11c2644c1220328a8773922ffb17733b7b47bc6081c55cb55",
	}
}
