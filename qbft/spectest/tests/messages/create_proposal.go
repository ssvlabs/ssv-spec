package messages

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() tests.SpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateProposal,
		Name:         "create proposal",
		Value:        [32]byte{1, 2, 3, 4},
		ExpectedRoot: "b28e854b54d18a9be6450ad72b6d6e3724ae83f74c185095b8b721c09771c2d1",
	}
}
