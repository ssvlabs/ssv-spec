package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidValCheck tests a quorum of commits received with an invalid value check
func InvalidValCheck() tests.SpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{}
	// No need to check as a commit depends on a proposal received which validates value
	return &tests.MsgProcessingSpecTest{
		Name:           "commit invalid val check",
		Pre:            pre,
		PostRoot:       "35d96d47901971b46055a010e136a97c71644d7d8bee368a8c1f29c149d3fa6c",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
