package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateProposal,
		Name:         "create proposal",
		Value:        [32]byte{1, 2, 3, 4},
		Round:        10,
		ExpectedRoot: "bef9eabbeaee0d8067760f17aabe154af09a6678fbe65c370a68421e3a0ecfd2",
	}
}
