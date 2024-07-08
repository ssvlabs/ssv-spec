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

// PostDecided tests a valid start duty before finished and after decided of another duty.
// Duties that have a preconsensus phase won't update the `currentRunningInstance`.
func PostDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	decidedRunner := func(r ssv.Runner, duty types.Duty) ssv.Runner {
		// baseStartNewDuty(r, duty) will override this state.
		// We set it here to correctly mimic the state of the runner after the duty is started.
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.CommitteeMember,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.DutySlot()))
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.DutySlot())
		return r
	}

	multiSpecTest := &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  decidedRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: postDecidedSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postDecidedSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  decidedRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: postDecidedAggregatorSC().Root(),
				PostDutyRunnerState:     postDecidedAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  decidedRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterDuty),
				Duty:                    testingutils.TestingAttesterDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: postDecidedAttesterSC().Root(),
				PostDutyRunnerState:     postDecidedAttesterSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
			{
				Name:                    "sync committee",
				Runner:                  decidedRunner(testingutils.CommitteeRunner(ks), testingutils.TestingSyncCommitteeDuty),
				Duty:                    testingutils.TestingSyncCommitteeDutyNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: postDecidedSyncCommitteeSC().Root(),
				PostDutyRunnerState:     postDecidedSyncCommitteeSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
			{
				Name:                    "attester and sync committee",
				Runner:                  decidedRunner(testingutils.CommitteeRunner(ks), testingutils.TestingAttesterAndSyncCommitteeDuties),
				Duty:                    testingutils.TestingAttesterAndSyncCommitteeDutiesNextEpoch,
				Threshold:               ks.Threshold,
				PostDutyRunnerStateRoot: postDecidedCommitteeSC().Root(),
				PostDutyRunnerState:     postDecidedCommitteeSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer (%s)", version.String()),
			Runner:                  decidedRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			Threshold:               ks.Threshold,
			PostDutyRunnerStateRoot: postDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     postDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:                  decidedRunner(testingutils.ProposerBlindedBlockRunner(ks), testingutils.TestingProposerDutyV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			Threshold:               ks.Threshold,
			PostDutyRunnerStateRoot: postDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postDecidedBlindedProposerSC(version).ExpectedState,
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
