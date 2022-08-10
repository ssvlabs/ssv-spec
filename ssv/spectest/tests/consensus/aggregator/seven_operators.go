package aggregator

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
	dr := testingutils.AggregatorRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[5], 5)),

		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestAggregatorConsensusDataByts, nil, nil),
		}), nil),

		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),

		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),
		testingutils.SSVMsgAggregator(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AggregatorMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
		}), nil),

		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[5], 5)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "aggregator 7 operator happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestAggregatorConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "940fa3c836253ea3fa3c43a2fe8e2db7e8afd57453ac4599036e41332be5e5b4",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusSelectionProofMsg(testingutils.Testing7SharesSet().Shares[1], 1),
			testingutils.PostConsensusAggregatorMsg(testingutils.Testing7SharesSet().Shares[1], 1),
		},
	}
}
