package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFutureDecided tests starting duty after a future decided
func PostFutureDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	futureDecide := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)

		futureDecidedInstance := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			50)
		futureDecidedInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, futureDecidedInstance)
		r.GetBaseRunner().QBFTController.Height = 50
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post future decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  futureDecide(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "4d90d1d592e40bf30f539f1619b549c30304e7033207ddec7b8621485f49e1f4",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  futureDecide(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "084bcbd0f7ae46e6f80430adbdfac088405e18af1cdc9180ed217506e95c5ed1",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  futureDecide(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "325d39e970c78c6842f0b0b68bc1790b329dc1436a5cfd9846443f3f1153d33e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  futureDecide(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "2d7fc45c36a11a6cadf4e4921e169d36b03930522fb07c325503b708aa9c57f6",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  futureDecide(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "1c05912a3eedf8c144f71dd0d950ca7f45cd13cf3850df53d408948c9eeee0ba",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
