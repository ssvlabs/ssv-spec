package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// OldSlotDuty tests starting duty from an old slot after we ran a duty from a higher slot.
// This shouldn't happen, and we should return an error.
func OldSlotDuty() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	futureDecide := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		duty.Slot += 50
		r.GetBaseRunner().State = ssv.NewRunnerState(3, &duty)
		return r
	}

	expectedError := "can't start duty: duty for slot 12 already passed. Current height is 62"

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty old slot",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  futureDecide(testingutils.SyncCommitteeContributionRunner(ks), testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "0c2bf5b2570aad7da85cfe5aa81361524a727265c22b41e72e7e9ff2f9b2c215",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "sync committee",
				Runner:                  futureDecide(testingutils.SyncCommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "89edc9d9c28654a0113c3003c1538aaae36fe14490992dbf057b9f5e5d492e33",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  futureDecide(testingutils.AggregatorRunner(ks), testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "892d0c2842a81d163c59e53302085f7bb2753cc4fcb46593bfeb569ec39e8928",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  futureDecide(testingutils.ProposerRunner(ks), *testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "fdaaa35d42c3001cd891209a44b921fa64320be238794e01633661e16c4f5e02",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "attester",
				Runner:                  futureDecide(testingutils.AttesterRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "ca53abb401eaae1154b075d5fc6ddca2da760c097fc30da8ee8e3abb94efb6d2",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
		},
	}
}
