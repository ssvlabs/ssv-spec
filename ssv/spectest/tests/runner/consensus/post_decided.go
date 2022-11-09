package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
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
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeContributionConsensusDataRoot, Source: testingutils.TestSyncCommitteeContributionConsensusDataByts}, ks, types.BNRoleSyncCommitteeContribution),
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestSyncCommitteeContributionConsensusDataRoot,
							Source: testingutils.TestSyncCommitteeContributionConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType)),
				PostDutyRunnerStateRoot: "2474d17fa64a0f23e0c5e760a1ed11ea16bca4af26499b8d85cc4284221092d0",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeConsensusDataRoot, Source: testingutils.TestSyncCommitteeConsensusDataByts}, ks, types.BNRoleSyncCommittee),
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestSyncCommitteeConsensusDataRoot,
							Source: testingutils.TestSyncCommitteeConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType)),
				PostDutyRunnerStateRoot: "a8aa102d5294d08e1e17a58d221a6d992e6c7bc6e884abf46bf6148912b6bd8e",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAggregatorConsensusDataRoot, Source: testingutils.TestAggregatorConsensusDataByts}, ks, types.BNRoleAggregator),
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestAggregatorConsensusDataRoot,
							Source: testingutils.TestAggregatorConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType)),
				PostDutyRunnerStateRoot: "9bc1eb96d7a5073b88f23c11922303f0b29848c8cd86ae1c46888c1819054b22",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestProposerConsensusDataRoot, Source: testingutils.TestProposerConsensusDataByts}, ks, types.BNRoleProposer),
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestProposerConsensusDataRoot,
							Source: testingutils.TestProposerConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType)),
				PostDutyRunnerStateRoot: "66a34673fa59a3bf6b3d001081d01d01cbacb75e5f7dc0f8f42b48b26ffedd77",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAttesterConsensusDataRoot, Source: testingutils.TestAttesterConsensusDataByts}, ks, types.BNRoleAttester),
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							Height: qbft.FirstHeight,
							Round:  qbft.FirstRound,
						}, &qbft.Data{
							Root:   testingutils.TestAttesterConsensusDataRoot,
							Source: testingutils.TestAttesterConsensusDataByts,
						}), nil, types.ConsensusCommitMsgType)),
				PostDutyRunnerStateRoot: "4ac261058e32f6bcd6a96e2fa52081c5ca5f830a7724d4bdfc6cc1a988c002ac",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PostConsensusAttestationMsg(testingutils.Testing4SharesSet().Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
