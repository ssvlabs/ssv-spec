package consensus

import (
	"crypto/sha256"
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
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Root:       sha256.Sum256(testingutils.TestSyncCommitteeContributionConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "ae7ee7e668eb90c4188316a8db8a4b36cd329bc859dd997cc5b9be60fb4e89e7",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
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
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Root:       sha256.Sum256(testingutils.TestSyncCommitteeConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "e2bcc9c4304b5e8409628039f900e2dac104ca75f7b2d12f1cd52b45b44eac75",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
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
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Root:       sha256.Sum256(testingutils.TestAggregatorConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "f81db6a51506ad7a4ccb972e84eabd7f45132ea825f06df7c65c07d8583fe006",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
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
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Root:       sha256.Sum256(testingutils.TestProposerConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "e2f2c71d4ffd08093bc57b937f5c2618973315f92835559bf3ea272a1c033517",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
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
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Root:       sha256.Sum256(testingutils.TestProposerBlindedBlockConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "5e982dd86d046bb26cee80362a4f7937653dae0377244bc2a4553d487942640a",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
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
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Root:       sha256.Sum256(testingutils.TestAttesterConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "14771f0b1b77704b58a9c62a9ce06566964747ddf60f3f9c280b356b31a5d126",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(
						testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ValidatorRegistrationMsgID,
							Root:       sha256.Sum256(testingutils.TestAttesterConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "34b8788e99ec86321d758e7a55a7e688bb91cc15d0a8fd581cf28f0839601fbd",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
		},
	}
}
