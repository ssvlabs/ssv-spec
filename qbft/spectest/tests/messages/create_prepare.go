package messages

import "github.com/bloxapp/ssv-spec/qbft/spectest/tests"

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreatePrepare,
		Name:         "create prepare",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "106517076940f814d71164a38301e5945e8280c05e6938f4eb6c3e0c37eaa513",
	}
}
