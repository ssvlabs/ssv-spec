package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongRandaoRoot tests processing randao msg with the wrong signing root
func WrongRandaoRoot() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentEpochMsg(ks.Shares[1], 1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao wrong root",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "9ffe08b74a87f4c5395b7a3978b7be6f45709cec0b924eac6fb5543bc621a1c4",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedError: "failed processing randao message: invalid randao message: wrong randao signing root",
	}
}
