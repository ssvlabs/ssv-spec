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
		PostDutyRunnerStateRoot: "0485e75de7b087b605e2121d505414d56f46a5368786b9429762ea271f593e2b",
		ExpectedError:           "unknown msg",
	}
}
