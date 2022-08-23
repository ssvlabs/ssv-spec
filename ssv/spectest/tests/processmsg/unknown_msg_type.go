package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownMsgType tests an unknown SSVMessage type
func UnknownMsgType() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msg := testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight))
	msg.MsgType = 100
	msgs := []*types.SSVMessage{
		msg,
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "unknown msg type",
		Runner:                  dr,
		Messages:                msgs,
		Duty:                    testingutils.TestingAttesterDuty,
		PostDutyRunnerStateRoot: "79338fd23b14bee97cb0356da5bf6c6c30ff517fe62fc4dd32fcb440d2725284",
		ExpectedError:           "unknown msg",
	}
}
