package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsgs tests a processing duplicate msgs
func DuplicateMsgs() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao duplicate msg",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "3f55ac00cd2b200bf7ad8784a38f60bc9f1389d110724a35bfed64aa29cde9c6",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
	}
}
