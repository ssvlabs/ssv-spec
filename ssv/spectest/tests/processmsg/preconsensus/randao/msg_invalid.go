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
		PostDutyRunnerStateRoot: "ca3d758a37f4448b654c844b2990ea8fe705920ee31b9732ed7bcb92ac7e5400",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedError: "failed processing randao message: invalid pre-consensus message: SignedPartialSignatureMessage invalid: message invalid: SigningRoot invalid",
	}
}
