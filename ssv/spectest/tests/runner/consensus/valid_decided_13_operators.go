package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided13Operators tests a valid decided value (13 operators)
func ValidDecided13Operators() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing13SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 13 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeContributionConsensusDataRoot, Source: testingutils.TestSyncCommitteeContributionConsensusDataByts}, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "e3075ce2f461f2136eec3fa17c0327bd27a389169b217497d630c50ad9d8166e",
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
				PostDutyRunnerStateRoot: "b3a3d97c4c6c3662218c5ddae5b068642cea42cf69f915975140bba7a33bbb9a",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAggregatorConsensusDataRoot, Source: testingutils.TestAggregatorConsensusDataByts}, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "5aa751317f1cdd4a64df6e2a977275b905d1fc0a3e1874ded835b7b0064081a4",
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
				PostDutyRunnerStateRoot: "2c804ce40938f0985e086c1f020296abc500d0aad1ed04898c2624b488876710",
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
				PostDutyRunnerStateRoot: "f76f90ec622f79984a7e6121efadadb8b2e73f4db87d7f847d8ffbccbabbaaab",
				OutputMessages: []*ssv.SignedPartialSignatures{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
