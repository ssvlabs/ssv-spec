package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PostFinish tests a valid commit msg after runner finished
func PostFinish() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		// *qbft.Instance cannot be shared between StoredInstances and RunningInstance
		// because JSON decoder creates a new *qbft.Instance for each of them
		newInstance := func() *qbft.Instance {
			return qbft.NewInstance(
				r.GetBaseRunner().QBFTController.GetConfig(),
				r.GetBaseRunner().Share,
				r.GetBaseRunner().QBFTController.Identifier,
				qbft.FirstHeight)
		}
		r.GetBaseRunner().State.RunningInstance = newInstance()
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, newInstance())
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		r.GetBaseRunner().State.Finished = true
		return r
	}

	err := "failed processing consensus message: no running duty"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: finishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						), nil),
				},
				PostDutyRunnerStateRoot: "5416fdd1e8e5ac26d01bce867630c8956b6b8444b1d92e0f2d2b1e23322dba24",
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
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						), nil),
				},
				PostDutyRunnerStateRoot: "781e504c4c8d952a4145fed6453f7709fccf632f597ec5bc3071cbd6aca8522b",
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
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						), nil),
				},
				PostDutyRunnerStateRoot: "dc23550fb81816c01ca495936ead9b98469d27b9f6f3a1edd871b11934724ad3",
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
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataByts,
						), nil),
				},
				PostDutyRunnerStateRoot: "b165432cb8ec536de81e16cc2de5ae621c787931530f26fee235a9f2bc368cfb",
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
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataByts,
						), nil),
				},
				PostDutyRunnerStateRoot: "0d55336a5a0f0f2624d0396b3b4f1f5509ff69fe0e14708d18c69c92190a6966",
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
						testingutils.TestingCommitMultiSignerMessageWithIdentifierAndFullData(
							[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						), nil),
				},
				PostDutyRunnerStateRoot: "c9b543054dd8a3a051ae68599171c4e7dc926f9ba0862d4d75aa90da248c5956",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           err,
			},
		},
	}
}
