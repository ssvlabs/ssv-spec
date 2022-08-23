package randao

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DutyFinished tests processing a randao msg when the duty is already completed
func DutyFinished() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.ProposerRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[3], 3)),

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

		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),

		testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsg(ks.Shares[4], 4)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "randao duty finished",
		Runner:                  dr,
		Duty:                    testingutils.TestingProposerDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "8e42807f6c48ddf194b200c017bf2673d19effe730a7bfd6965e790062b41410",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
			testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1),
		},
		ExpectedError: "Message invalid: no running duty",
	}
}
