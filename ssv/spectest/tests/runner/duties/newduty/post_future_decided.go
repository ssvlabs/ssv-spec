package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostFutureDecided tests starting duty after a future decided
// This can happen if we receive a future decided message from the network and we are behind.
func PostFutureDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	futureDecide := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.DutySlot()))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)

		futureDecidedInstance := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.DutySlot()+50))
		futureDecidedInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, futureDecidedInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.DutySlot() + 50)
		return r
	}

	expectedError := "can't start duty: duty for slot 12 already passed. Current height is 62"

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post future decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  futureDecide(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "0c2bf5b2570aad7da85cfe5aa81361524a727265c22b41e72e7e9ff2f9b2c215",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  futureDecide(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "892d0c2842a81d163c59e53302085f7bb2753cc4fcb46593bfeb569ec39e8928",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  futureDecide(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: "fdaaa35d42c3001cd891209a44b921fa64320be238794e01633661e16c4f5e02",
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: fmt.Sprintf("can't start duty: duty for slot %v already passed. Current height is %v",
					testingutils.TestingDutySlotV(spec.DataVersionDeneb),
					testingutils.TestingDutySlotV(spec.DataVersionDeneb)+50),
			},
			{
				Name:           "attester",
				Runner:         futureDecide(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty),
				Duty:           testingutils.TestingAttesterDuty,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:           "sync committee",
				Runner:         futureDecide(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:           testingutils.TestingSyncCommitteeDuty,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:           "attester and sync committee",
				Runner:         futureDecide(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties),
				Duty:           testingutils.TestingAttesterAndSyncCommitteeDuties,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		},
	}
}
