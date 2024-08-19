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
		ExpectedRoot: "6f6ad56759e8e7dc35135b6780f2aa3b0f7e6a7efe34d74dc80fd819c6cfbfe5",
	}
}
