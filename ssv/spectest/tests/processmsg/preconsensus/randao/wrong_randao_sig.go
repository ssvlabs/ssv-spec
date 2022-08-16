package randao

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongRandaoSigner tests an invalid randao sig over root
func WrongRandaoSigner() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msg := testingutils.PreConsensusRandaoMsg(ks.Shares[1], 2)
	msg.Signer = 1 // it signed the randao root with operator #2 and now we change the msg signer to 1 so it won't fail valdiation but the randao sig will
	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, msg),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao wrong randao signer",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "9ffe08b74a87f4c5395b7a3978b7be6f45709cec0b924eac6fb5543bc621a1c4",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
		},
		ExpectedError: "failed processing randao message: invalid randao message: could not verify Beacon partial Signature: wrong signature",
	}
}
