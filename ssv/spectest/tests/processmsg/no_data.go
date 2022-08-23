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
		PostDutyRunnerStateRoot: "79338fd23b14bee97cb0356da5bf6c6c30ff517fe62fc4dd32fcb440d2725284",
		ExpectedError:           "Message invalid: msg data is invalid",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
	}
}
