package attester

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ThirteenOperators tests a full attestation duty execution with 13 operators
func ThirteenOperators() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing13SharesSet()
	dr := testingutils.AttesterRunner7Operators(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
		}), nil),

		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[6], 6, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[7], 7, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[8], 8, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[9], 9, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),

		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[6], 6, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[7], 7, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[8], 8, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[9], 9, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
		}), nil),

		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[4], 4, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[5], 5, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[6], 6, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[7], 7, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[8], 8, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[9], 9, qbft.FirstHeight)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "attester 13 operators happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestAttesterConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "29f2b38df73cd832a3c032dff3ba8b830bdf21a396c30ad2c5eb60213302091e",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PostConsensusAttestationMsg(testingutils.Testing13SharesSet().Shares[1], 1, qbft.FirstHeight),
		},
	}
}
