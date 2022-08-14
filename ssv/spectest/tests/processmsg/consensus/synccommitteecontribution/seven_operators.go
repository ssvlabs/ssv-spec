package synccommitteecontribution

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
	dr := testingutils.SyncCommitteeContributionRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], 2)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], 3)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], 4)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[5], 5)),

		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts, nil, nil),
		}), nil),

		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.PrepareDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),

		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[2], 2, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),
		testingutils.SSVMsgSyncCommitteeContribution(testingutils.SignQBFTMsg(ks.Shares[5], 5, &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.SyncCommitteeContributionMsgID,
			Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
		}), nil),

		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[4], 4, ks)),
		testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[5], 5, ks)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "sync committee contribution 7 operators happy flow",
		Runner:                  dr,
		Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "6a42fc9384226a06835027bb71192675898b2e8aaf23ff5a73d7228b00900382",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PreConsensusContributionProofMsg(testingutils.Testing7SharesSet().Shares[1], 1),
			testingutils.PostConsensusSyncCommitteeContributionMsg(testingutils.Testing7SharesSet().Shares[1], 1, testingutils.Testing7SharesSet()),
		},
	}
}
