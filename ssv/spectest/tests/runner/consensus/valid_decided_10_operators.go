package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
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
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeContributionConsensusDataRoot, Source: testingutils.TestSyncCommitteeContributionConsensusDataByts}, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "401e6d6e14308cd0c217dc2d84762b9fa894944a2b3da4cde86ebcc1e0674edc",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeConsensusDataRoot, Source: testingutils.TestSyncCommitteeConsensusDataByts}, ks, types.BNRoleSyncCommittee),
				PostDutyRunnerStateRoot: "d76948a3d6e95b14a17f81cc7a1b7f463b470a39990d42789fc18e5ff70c026b",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAggregatorConsensusDataRoot, Source: testingutils.TestAggregatorConsensusDataByts}, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "92727ebc9004a40a70123f588236fdf305110aaa07d6ce16fc1b30802ec47dc7",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestProposerConsensusDataRoot, Source: testingutils.TestProposerConsensusDataByts}, ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "c0ed7eab54f019b82c0fb8dc59be702b8d45710a681b919305b04b11ad0916f1",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAttesterConsensusDataRoot, Source: testingutils.TestAttesterConsensusDataByts}, ks, types.BNRoleAttester),
				PostDutyRunnerStateRoot: "8a19f42e9e40eaf1ad95ef6b2b40e0098ae284ad0485b9a4639daf53f14a22d6",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
