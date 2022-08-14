package attester

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidConsensusMsg tests an invalid consensus msg processing
func InvalidConsensusMsg() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     100,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
		}), nil),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "attester invalid consensus msg",
		Runner:                  dr,
		Duty:                    testingutils.TestAttesterConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "79338fd23b14bee97cb0356da5bf6c6c30ff517fe62fc4dd32fcb440d2725284",
		OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
		ExpectedError:           "failed processing consensus message: failed to process consensus msg: instance not found",
	}
}
