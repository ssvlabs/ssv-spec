package processmsg

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoData tests a SSVMessage with no data
func NoData() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester),
			Data:    nil,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "ssv msg no data",
		Runner:                  dr,
		Messages:                msgs,
		Duty:                    testingutils.TestingAttesterDuty,
		PostDutyRunnerStateRoot: "0485e75de7b087b605e2121d505414d56f46a5368786b9429762ea271f593e2b",
		ExpectedError:           "Message invalid: msg data is invalid",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
	}
}
