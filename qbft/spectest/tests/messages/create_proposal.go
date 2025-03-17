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
		ExpectedRoot: "82802b9530012981e50ab2ea72ea2e914e584106a8c5795a2dc0a2ce494cafa3",
	}
}
