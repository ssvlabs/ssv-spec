package processmsg

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownRunner tests an SSVMessage with msg ID not matching any runner
func UnknownRunner() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)
	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingValidatorPubKey[:], 100),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "no runner for msg id",
		Runner:                  dr,
		Messages:                msgs,
		Duty:                    testingutils.TestingAttesterDuty,
		PostDutyRunnerStateRoot: "79338fd23b14bee97cb0356da5bf6c6c30ff517fe62fc4dd32fcb440d2725284",
		ExpectedError:           "could not get duty runner for msg ID",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
	}
}
