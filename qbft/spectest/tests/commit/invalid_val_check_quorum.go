package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidValCheck tests a quorum of commits received with an invalid value check
func InvalidValCheck() tests.SpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*types.SignedSSVMessage{}
	// No need to check as a commit depends on a proposal received which validates value
	return &tests.MsgProcessingSpecTest{
		Name:           "commit invalid val check",
		Pre:            pre,
		PostRoot:       "620ad2417e47411537db8df9d4a072327e3c3efc391c3162867f30d5bf9af52c",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
