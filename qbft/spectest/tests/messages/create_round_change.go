package messages

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:      tests.CreateRoundChange,
		Name:            "create round change",
		Value:           testingutils.TestingQBFTFullData,
		ExpectedSSZRoot: "cb741776552a93ba482e6eac9821fe8a83a478d96aba36249066bd635e617cd5",
	}
}
