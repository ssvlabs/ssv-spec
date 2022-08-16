package processmsg

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningDuty tests an SSVMessage with msg ID matching a runner but for which no duty started
func NoRunningDuty() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)
	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleProposer),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "no running duty",
		Runner:                  dr,
		Messages:                msgs,
		Duty:                    testingutils.TestingAttesterDuty,
		PostDutyRunnerStateRoot: "79338fd23b14bee97cb0356da5bf6c6c30ff517fe62fc4dd32fcb440d2725284",
		ExpectedError:           "Message invalid: no running duty",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
	}
}
