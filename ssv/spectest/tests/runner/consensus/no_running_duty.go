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
				Duty: testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "8856db8a55ac7054ac8c0d778c417779d83a81fd0f3eaea377435782258d2d0b",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "sync committee",
				Runner: startInstance(
					testingutils.SyncCommitteeRunner(ks),
					testingutils.TestSyncCommitteeConsensusDataByts,
				),
				Duty: testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestSyncCommitteeConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "6ff3b0bed9c4119edc612a81218369088b7c9021bf978bf7ce25b7f751fb6a14",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "aggregator",
				Runner: startInstance(
					testingutils.AggregatorRunner(ks),
					testingutils.TestAggregatorConsensusDataByts,
				),
				Duty: testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAggregatorConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "2521f4305f88dc09186af4c8dec213e678e7f45561d469ad462a9a4d4fc38a3e",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "proposer",
				Runner: startInstance(
					testingutils.ProposerRunner(ks),
					testingutils.TestProposerConsensusDataByts,
				),
				Duty: testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestProposerConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "8c705379bc2c68a36f7358530e46e76ad7477a6eca8e0db386de12aa577b9baf",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name: "attester",
				Runner: startInstance(
					testingutils.AttesterRunner(ks),
					testingutils.TestAttesterConsensusDataByts,
				),
				Duty: testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "7855ef3316a289373336d44d3f8ae43ea16bb5a6f04143e9b8ab060757c65b12",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ValidatorRegistrationMsgID,
							Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
						}), nil),
				},
				PostDutyRunnerStateRoot: "f5ae91f93c1352e8e94ef1564e6cf298c4244120d843bf1733417770ff6c65af",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
		},
	}
}
