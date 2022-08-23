package proposer

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SevenOperators tests a full valcheck + post valcheck + duty sig reconstruction flow for 7 operators
func SevenOperators() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing7SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[5], 5)),

		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestProposerConsensusDataByts, nil, nil),
		}), nil),

		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),

		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),
		testingutils.SSVMsgProposer(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.ProposerMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
		}), nil),

		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[5], 5)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "proposer 7 operators happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestProposerConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "f78a60ba301526f68d47989b18af08282383ac4db1880ca82748c0b9cc32b9d0",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(testingutils.Testing7SharesSet().Shares[1], 1),
			testingutils.PostConsensusProposerMsg(testingutils.Testing7SharesSet().Shares[1], 1),
		},
	}
}
