package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningDuty tests a valid proposal msg before duty starts
func NoRunningDuty() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	startInstance := func(r ssv.Runner, value []byte) ssv.Runner {
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight))

		return r
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: startInstance(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestSyncCommitteeContributionConsensusDataByts,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "5bdca31d56869d967df58ebedd89e130071a50dd92a482c79b01a342251b41e0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "sync committee",
				Runner: startInstance(
					testingutils.SyncCommitteeRunner(ks),
					testingutils.TestSyncCommitteeConsensusDataByts,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "b99077b644b160fd35d1a138a68ea9abf21990d74df22a1284e3f5cf2b527350",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "aggregator",
				Runner: startInstance(
					testingutils.AggregatorRunner(ks),
					testingutils.TestAggregatorConsensusDataByts,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "8b1322320222b9894052c6c17d621535a37ece0d5ae7f324096a23c4c50a3b02",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "proposer",
				Runner: startInstance(
					testingutils.ProposerRunner(ks),
					testingutils.TestProposerConsensusDataByts,
				),
				Duty: &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "483814963f6daeb4a730d1adad7d0524eb4edf224f98e201b4261786a8f88d71",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "proposer (blinded block)",
				Runner: startInstance(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestProposerBlindedBlockConsensusDataByts,
				),
				Duty: &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "dbbd72ea8094fb8be921f44f383fbd8880c05abf84d960a37a4fb28bbd743b18",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "attester",
				Runner: startInstance(
					testingutils.AttesterRunner(ks),
					testingutils.TestAttesterConsensusDataByts,
				),
				Duty: &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "7508cb14dd0d5833d218141b74b6fd8c512a54e99143c7f91d754ea3f08f1bd6",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ValidatorRegistrationMsgID,
							testingutils.TestAttesterConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
		},
	}
}
