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
		PostDutyRunnerStateRoot: "0485e75de7b087b605e2121d505414d56f46a5368786b9429762ea271f593e2b",
		ExpectedError:           "could not get duty runner for msg ID",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
	}
}
