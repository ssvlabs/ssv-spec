package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// ValidDecidedConsensusMsg tests an valid decided dutyexe msg
func ValidDecidedConsensusMsg() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	proposalByts, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier[:],
		Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
	}).Encode()
	decidedByts, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{
			testingutils.Testing4SharesSet().Shares[1],
			testingutils.Testing4SharesSet().Shares[2],
			testingutils.Testing4SharesSet().Shares[3],
		},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: identifier[:],
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}).Encode()
	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   identifier,
			Data:    proposalByts,
		},
		{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   identifier,
			Data:    decidedByts,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "valid decided dutyexe msg",
		Runner:                  dr,
		Messages:                msgs,
		Duty:                    testingutils.TestingAttesterDuty,
		PostDutyRunnerStateRoot: "8cee4711a21f23d06ef833a2742248a4604144e5cedacf4a83635086054d99cd",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PostConsensusAttestationMsg(testingutils.Testing4SharesSet().Shares[1], 1, qbft.FirstHeight),
		},
	}
}
