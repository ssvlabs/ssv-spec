package newduty

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ConsensusNotStarted tests starting duty after prev already started but for some duties' consensus didn't start because pre-consensus didnt get quorum (different duties will enable starting a new duty)
func ConsensusNotStarted() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	startRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty consensus not started",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  startRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "c78b42fb7db3a7683a78c2dd6a672d8d007c0997e580440c346b02d6f60f40b3",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  startRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDutyNextEpoch,
				PostDutyRunnerStateRoot: "6309028d3972d81c7ed985ed9e20b2c45779fb1c254bd36b38e66e381572d953",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "f5fd1090974190cecfecd5ffbc1f55f8b17c9c1b8f6c4e2888412517c8fb8e73",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "b2ae77c8491e702c0fd87114560cb447808406fa43afc72c163c533227d771c2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  startRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDutyNextEpoch,
				PostDutyRunnerStateRoot: "198b4b184304c99c41b4c161bf33c1427a727f520ef946e29f4880c11646b1a3",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "voluntary exit",
				Runner:                  startRunner(testingutils.VoluntaryExitRunner(ks), &testingutils.TestingVoluntaryExitDuty),
				Duty:                    &testingutils.TestingVoluntaryExitDutyNextEpoch,
				PostDutyRunnerStateRoot: "6f6d918e15ebc7b84cb77e2d603019d1cbfb6d7293daddd48780da47c14e53ce",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusVoluntaryExitNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "validator registration",
				Runner:                  startRunner(testingutils.ValidatorRegistrationRunner(ks), &testingutils.TestingValidatorRegistrationDuty),
				Duty:                    &testingutils.TestingValidatorRegistrationDutyNextEpoch,
				PostDutyRunnerStateRoot: "6f6d918e15ebc7b84cb77e2d603019d1cbfb6d7293daddd48780da47c14e53ce",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationNextEpochMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
