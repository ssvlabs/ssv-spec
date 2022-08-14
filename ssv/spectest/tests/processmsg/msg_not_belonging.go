package processmsg

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgNotBelonging tests an SSVMessage ID that doesn't belong to the validator
func MsgNotBelonging() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		{
			MsgType: 100,
			MsgID:   types.NewMsgID(testingutils.TestingWrongValidatorPubKey[:], types.BNRoleAttester),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "ssv msg wrong pubkey in msg id",
		Runner:                  dr,
		Messages:                msgs,
		Duty:                    testingutils.TestingAttesterDuty,
		PostDutyRunnerStateRoot: "79338fd23b14bee97cb0356da5bf6c6c30ff517fe62fc4dd32fcb440d2725284",
		ExpectedError:           "Message invalid: msg ID doesn't match validator ID",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
	}
}
