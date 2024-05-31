package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidValCheck tests a quorum of commits received with an invalid value check
func InvalidValCheck() tests.SpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*types.SignedSSVMessage{}
	// No need to check as a commit depends on a proposal received which validates value
	return &tests.MsgProcessingSpecTest{
		Name:           "commit invalid val check",
		Pre:            pre,
		PostRoot:       "01489f7af13579b66ce3da156d4d10208c85a10365380f04e7b8d82d0a9679ce",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
