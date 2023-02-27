package consensus

import (
	"crypto/sha256"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests a valid commit msg after returned decided already
func PostDecided() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Root:       sha256.Sum256(testingutils.TestSyncCommitteeContributionConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "4dfe9759ce5e43e7b131d7557bb9d03b0bab59b30ce07e8cc6876cdcf2281448",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Root:       sha256.Sum256(testingutils.TestSyncCommitteeConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "581bbfc385e1839cf3d894eec1a8a713a628f97520552d1f3da5de3aec15c42e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Root:       sha256.Sum256(testingutils.TestAggregatorConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "274174f3fa6f374a528e686f3d7152cc16bfe4405a62b74db4b6ed8aa99bcfd6",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusData, ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Root:       sha256.Sum256(testingutils.TestProposerConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "254ae86a32374da4b5fe4e4c61d6f20a00eecbe1e8cee0c22b1aec561c988cd8",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerBlindedBlockConsensusData, ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Root:       sha256.Sum256(testingutils.TestProposerBlindedBlockConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "c087a558b06707584a160e0bddc0bc0bdd265d049400af33559c3cda7436207c",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Root:       sha256.Sum256(testingutils.TestAttesterConsensusDataByts),
						}), nil)),
				PostDutyRunnerStateRoot: "7c72a4080295cdc2d7289366747d44327fa263e1ef901f8c1e9c00dcebc37865",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
