package messages

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateProposal,
		Name:         "create proposal",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "dfb0a692281b916b1d037df44f5c742f13ac3ee207ea0082cc3ca2afff34e178",
	}
}
