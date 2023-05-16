package messages

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:      tests.CreateProposal,
		Name:            "create proposal",
		Value:           []byte{1, 2, 3, 4},
		ExpectedSSZRoot: "443161d9ea4f2e4abd4c0545aeb4aa99bf40e8a22f59d9917c2dcb7f2f04c9f1",
	}
}
