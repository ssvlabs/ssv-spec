package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateRoundChange,
		Name:         "create round change",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "d31ba4f08c790d96b6951cf632c51923ce3e77fbc896336b1def371476b13743",
	}
}
