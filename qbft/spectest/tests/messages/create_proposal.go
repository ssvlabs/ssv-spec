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
		ExpectedRoot: "43c23219aaf744537a2e8b3896937a2e9aa24a8eaf8fcabf8ec7376f76669f3c",
	}
}
