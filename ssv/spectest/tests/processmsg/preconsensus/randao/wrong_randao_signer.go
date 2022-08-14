package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownRandaoSigner tests unknown PostConsensusMessage root signer
func UnknownRandaoSigner() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], testingutils.Testing7SharesSet().Shares[5], 1, 5)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao unknown randao signer",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "9ffe08b74a87f4c5395b7a3978b7be6f45709cec0b924eac6fb5543bc621a1c4",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedError: "failed processing randao message: invalid randao message: could not verify Beacon partial Signature: Beacon partial Signature Signer not found",
	}
}
