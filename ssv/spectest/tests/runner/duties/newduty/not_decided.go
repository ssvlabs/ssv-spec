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

// NotDecided tests starting duty before finished or decided
func NotDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	startRunner := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.DutySlot()))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.DutySlot())
		return r
	}

	multiSpecTest := &MultiStartNewRunnerDutySpecTest{
		Name: "new duty not decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  startRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: notDecidedSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     notDecidedSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: notDecidedAggregatorSC().Root(),
				PostDutyRunnerState:     notDecidedAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:           "attester",
				Runner:         startRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty),
				Duty:           testingutils.TestingAttesterDutyNextEpoch,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:           "sync committee",
				Runner:         startRunner(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:           testingutils.TestingSyncCommitteeDutyNextEpoch,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
			{
				Name:           "attester and sync committee",
				Runner:         startRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties),
				Duty:           testingutils.TestingAttesterAndSyncCommitteeDutiesNextEpoch,
				Threshold:      ks.Threshold,
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer (%s)", version.String()),
			Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			Threshold:               ks.Threshold,
			PostDutyRunnerStateRoot: notDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     notDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:                  startRunner(testingutils.ProposerBlindedBlockRunner(ks), testingutils.TestingProposerDutyV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			Threshold:               ks.Threshold,
			PostDutyRunnerStateRoot: notDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     notDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*StartNewRunnerDutySpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
