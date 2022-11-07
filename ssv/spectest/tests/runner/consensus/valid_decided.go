package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided tests a valid decided value
func ValidDecided() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestSyncCommitteeContributionConsensusDataRoot, Source: testingutils.TestSyncCommitteeContributionConsensusDataByts}, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "3d7bc2f3ce07747480eda6284334c684e9436e9c01b60b696d8257e829e365a4",
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
				PostDutyRunnerStateRoot: "3e927b2683a5e49730d15ac4dadb1dc307dcedab346345a9bdbf8822f17c0485",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAggregatorConsensusDataRoot, Source: testingutils.TestAggregatorConsensusDataByts}, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "3c7a0ab37d5bd1e0a5653c723dfd33acdd31f28766200710a1379370bffb19fb",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestProposerConsensusDataRoot, Source: testingutils.TestProposerConsensusDataByts}, ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "eda904f00bd990a997c1797bb86edcbd24bcbc0ab5762be2b2a125a77400b125",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PreConsensusRandaoMsg(testingutils.Testing4SharesSet().Shares[1], 1),
					testingutils.PostConsensusProposerMsg(testingutils.Testing4SharesSet().Shares[1], 1),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    testingutils.TestingAttesterDuty,
				Messages:                testingutils.SSVDecidingMsgs(&qbft.Data{Root: testingutils.TestAttesterConsensusDataRoot, Source: testingutils.TestAttesterConsensusDataByts}, ks, types.BNRoleAttester),
				PostDutyRunnerStateRoot: "d4c87cb532bc48908ddadad078c86e75f34d40002356cce3598cf4c569a0caf8",
				OutputMessages: []*ssv.SignedPartialSignature{
					testingutils.PostConsensusAttestationMsg(testingutils.Testing4SharesSet().Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
