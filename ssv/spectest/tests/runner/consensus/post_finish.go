package consensus

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		r.GetBaseRunner().State.Finished = true
		return r
	}

	err := "failed processing consensus message: could not process msg: invalid signed message: did not receive proposal for this round"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeContributionConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "d800b0a2047ef1caa04a7a97d803558f5f5e3c6efb5c6f0e28c94023d3676609",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "sync committee",
				Runner: finishRunner(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestSyncCommitteeConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "09544309e2adc09cac92cd0d5a0c6b9c894b9716989a2a365f46067d95b02c90",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestAggregatorConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "c05cba9bb0738aba2e4e2307f37315f82d28d872ada18e8bc5889ec32a7b8e98",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDuty),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestProposerConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "5735a6a3652729ec867d2ff943d529ad496becc88c4e3a680bcdc4376984e365",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: finishRunner(testingutils.ProposerBlindedBlockRunner(ks), testingutils.TestingProposerDuty),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestProposerBlindedBlockConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "3f667fcc28d662918d3b9eb47d257af09cfb5f25bfabc51ee5bbf8784d169fea",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "attester",
				Runner: finishRunner(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "b21d79156d14ad5c4e7e97b79dbdf3923cd19637cc9783baa8e113b6a82f05b5",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
		},
	}
}
