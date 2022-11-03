package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided7Operators tests a valid decided value (7 operators)
func ValidDecided7Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing7SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 7 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeContributionConsensusDataRoot, Source: testingutils.TestSyncCommitteeContributionConsensusDataByts}, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "bd8f45e4225e944488796ce8312293240fc7eec18651601129a2757f6c4010a0",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeConsensusDataRoot, Source: testingutils.TestSyncCommitteeConsensusDataByts}, ks, types.BNRoleSyncCommittee),
				PostDutyRunnerStateRoot: "23db669eb2d48842aa6f29be23f2d063c5c835a5e6eb22b61ffa7b043c191132",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAggregatorConsensusDataRoot, Source: testingutils.TestAggregatorConsensusDataByts}, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "6ad7c7b35d0dbaf0853a5a90db3459f3114708c10599ba08b2b057a21c0d2682",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestProposerConsensusDataRoot, Source: testingutils.TestProposerConsensusDataByts}, ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "d463e92e34f2b1cc2b4e85963288de4ddb92408ad552f5385aafc93b8fabf2b8",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAttesterConsensusDataRoot, Source: testingutils.TestAttesterConsensusDataByts}, ks, types.BNRoleAttester),
				PostDutyRunnerStateRoot: "0e2ae2ece4b027e1ac280faeda2ebf9d739ab057ab36cfb9cb52f3988c6519fc",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
