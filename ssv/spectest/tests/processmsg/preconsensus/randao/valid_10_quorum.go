package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Valid10Quorum tests a valid quorum (for 10 operators) of partial randao sig msg
func Valid10Quorum() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing10SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[5], 5)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[6], 6)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[7], 7)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao valid 10 quorum",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "5a3273032daa9d384d0be77fa1e37de1ca55092d58fa7ca9391e840deb42b9a7",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
	}
}
