package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidValCheck tests a quorum of commits received with an invalid value check
func InvalidValCheck() tests.SpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*types.SignedSSVMessage{}
	// No need to check as a commit depends on a proposal received which validates value
	return tests.NewMsgProcessingSpecTest(
		"commit invalid val check",
		testdoc.CommitTestInvalidValCheckDoc,
		pre,
		"",
		nil,
		msgs,
		nil,
		"",
		nil,
	)
}
