package consensus

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
)

// FutureDecided tests a running instance at FirstHeight, then processing a decided msg from height 10 and returning decided but doesn't move to post consensus as it's not the same instance decided
func FutureDecided() *tests.MultiMsgProcessingSpecTest {
	panic("implement")
}
