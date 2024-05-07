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
		ExpectedRoot: "a52caeaa01af9479834aeea962392d6df6ad373565f309700812c4062edb7840",
	}
}
