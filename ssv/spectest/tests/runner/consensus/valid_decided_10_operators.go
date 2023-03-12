package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided10Operators tests a valid decided value (10 operators)
func ValidDecided10Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing10SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 10 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "69ab7e6235e2d1d9c6425289f8bbfdefac86f06565287aa7d8a55b70063bb5bf",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
				PostDutyRunnerStateRoot: "073b633699bf5becffb156d95ae1c4bd245915f6f00dc483c746b9cc9d70e4e0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    &testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "69b8d4d30bba4722203d449f1f5ee405937511f770c94f97e26ae6d9b00d31bf",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    &testingutils.TestingProposerDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusData, ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "3df1cfe5434fbf28fe803fb986898830d536f60447a1d68a8ae36ad7d9b4a905",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer (blinded block)",
				Runner:                  testingutils.ProposerBlindedBlockRunner(ks),
				Duty:                    &testingutils.TestingProposerDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestProposerBlindedBlockConsensusData, ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "be3be758ef0829c238c76db01b6a579e52fc425f6cf5cda9eeeec1d2ccc01862",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    &testingutils.TestingAttesterDuty,
				Messages:                testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
				PostDutyRunnerStateRoot: "6fe16838e5e18f1b53870e2bb3431d8194a63cf2a7d5c483990da2fbe67c10c4",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
