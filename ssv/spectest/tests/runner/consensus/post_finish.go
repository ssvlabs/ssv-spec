package consensus

import (
	"crypto/sha256"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeContributionMsgID,
							Root:       sha256.Sum256(testingutils.TestSyncCommitteeContributionConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "dbc5a425e55618d83b245a71e1858419f8cd6d4b1b1022b68403da001c950601",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "sync committee",
				Runner: finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.SyncCommitteeMsgID,
							Root:       sha256.Sum256(testingutils.TestSyncCommitteeConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "47bbee3eb9895a8e3a96dde11075ed8ceb879f5548259ad06684f9a7c66034dc",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "aggregator",
				Runner: finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AggregatorMsgID,
							Root:       sha256.Sum256(testingutils.TestAggregatorConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "4a13cd997588052e33c46ecefee5621665d78dc4ffce1894a31cccf50124c2bc",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "proposer",
				Runner: finishRunner(testingutils.ProposerRunner(ks), &testingutils.TestingProposerDuty),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Root:       sha256.Sum256(testingutils.TestProposerConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "11a906a8bd9216788e2789bea9d935149ed232b41026888291a97a827d32ed67",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: finishRunner(testingutils.ProposerBlindedBlockRunner(ks), &testingutils.TestingProposerDuty),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.ProposerMsgID,
							Root:       sha256.Sum256(testingutils.TestProposerBlindedBlockConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "fc6c34a1ed8638d25426f5c6989e203c0e32016527f014d36e06dff79cb6a7ce",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
			{
				Name:   "attester",
				Runner: finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: testingutils.AttesterMsgID,
							Root:       sha256.Sum256(testingutils.TestAttesterConsensusDataByts),
						}), nil),
				},
				PostDutyRunnerStateRoot: "1978eeab8c871d7ee680f35dac401df80899efca5b6623f2ae1f6d1d64d05426",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
		},
	}
}
