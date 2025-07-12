package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create proposal",
		"Test creating a proposal message when not previously prepared.",
		[32]byte{1, 2, 3, 4},
		nil,
		10,
		nil,
		nil,
		tests.CreateProposal,
		"43c23219aaf744537a2e8b3896937a2e9aa24a8eaf8fcabf8ec7376f76669f3c",
		nil,
		"",
	)
}
