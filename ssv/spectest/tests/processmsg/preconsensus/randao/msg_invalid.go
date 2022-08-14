package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgInvalid tests an invalid partial randao sig msg (fails msg.Validate()) that fails to process
func MsgInvalid() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msg := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)
	msg.Message.Messages[0].SigningRoot = nil
	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, msg),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao invalid msg",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "9ffe08b74a87f4c5395b7a3978b7be6f45709cec0b924eac6fb5543bc621a1c4",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedError: "failed processing randao message: invalid randao message: SignedPartialSignatureMessage invalid: message invalid: SigningRoot invalid",
	}
}
