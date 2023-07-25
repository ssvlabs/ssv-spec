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
// This can happen if we receive a future decided message from the network and we are behind.
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
			qbft.Height(duty.Slot))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)

		futureDecidedInstance := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot+50))
		futureDecidedInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, futureDecidedInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot + 50)
		return r
	}

	expectedError := "can't start new duty runner instance for duty: could not start new QBFT instance: attempting to start an instance with a past height"

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post future decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  futureDecide(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "b4713150d68451ce58ffb8399509a7ebb674d80c88885898e364acac9dee287c",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  futureDecide(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "046ff8dbd9cbdf54cd8ce1356230f30deff71e47283275d01817e2d317b84194",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  futureDecide(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "c29585575507df93cbdbb0a6afcdeb541101c02fbb6aee346b5e48b29fcf9107",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  futureDecide(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "6d54dde37b0658f7d4447302e980ed505b063b23a87517014eada4646e126ffc",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  futureDecide(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "02bfc445dc650412518fb44d8a3a191ae8e180bcc4c2a4ac56e37485f4737d8f",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
		},
	}
}
