package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidDecided10Operators tests a valid decided value (10 operators)
func ValidDecided10Operators() tests.SpecTest {
	ks := testingutils.Testing10SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 10 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
				PostDutyRunnerStateRoot: "8a5d27764a97e43a36b87bd125291280390c7d5aef482255e3df4bb30a7989a1",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee),
				PostDutyRunnerStateRoot: "85441d6d9f2433e56436c7d0e8d3511fb9cc7536aba117c0d098ab0ce4362856",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    &testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator),
				PostDutyRunnerStateRoot: "de4c8388f01aa564d4649a64836f6f269700b5b7075132dba4be0364069aac1d",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(spec.DataVersionBellatrix), ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "e47546fe2de9b026280335577bc92d742547db1bd6937fa96116b74e552778ed",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
				},
			},
			{
				Name:                    "proposer (blinded block)",
				Runner:                  testingutils.ProposerBlindedBlockRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionBellatrix), ks, types.BNRoleProposer),
				PostDutyRunnerStateRoot: "a7c7cdbc1e2171a480e94f521402f5ba0918331b2d9eb90178aa68ea41b38907",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    &testingutils.TestingAttesterDuty,
				Messages:                testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester),
				PostDutyRunnerStateRoot: "ddb253e278c7fdce5bb25855dae4c15415ad5254963707a60c1adb540dd739a5",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
