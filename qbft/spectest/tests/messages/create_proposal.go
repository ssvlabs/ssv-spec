package messages

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType:   tests.CreateProposal,
		Name:         "create proposal",
		Value:        []byte{1, 2, 3, 4},
		ExpectedRoot: "d83735f8a99fabde910e4dff198163811be43ffbb1ffc62a80c75c44618270c1",
	}
}
