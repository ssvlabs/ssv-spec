package aggregator

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a full valcheck + post valcheck + duty sig reconstruction flow
func HappyFlow() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AggregatorRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),

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

		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "aggregator happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestAggregatorConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "35e59dc18886bdce93a2a82c05b9fd383f52b8f365f02d1938c19f3c6b0aa011",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusSelectionProofMsg(testingutils.Testing4SharesSet().Shares[1], 1),
			testingutils.PostConsensusAggregatorMsg(testingutils.Testing4SharesSet().Shares[1], 1),
		},
	}
}
