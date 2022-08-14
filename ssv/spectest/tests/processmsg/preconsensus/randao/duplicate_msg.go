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
		PostDutyRunnerStateRoot: "bf0fc94c6f5e64c39cd7f83e62e28db3fc15a4c8fa3e85f775062a14679bde87",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
	}
}
