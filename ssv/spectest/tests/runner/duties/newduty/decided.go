package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests a valid start duty before finished and after decided
func PostDecided() tests.SpecTest {
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
		return r
	}

	multiSpecTest := &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  finishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "cbb0ce082e86ba5526eb33a8b4c7f6f1e59d520f7b304beb894631e7d86a4765",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "0e6005bad8a43f831437f40a385926f2707cfd151ce239d995fc1c96b2a3c4e8",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "3b44e308331e22a88023ab437cf15ff64a6fa000be46703d27642529fb0acfec",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "326a86dd6ad7016e640f9f3418499552c8fac2691029a3812370017924c0bb0c",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer (%s)", version.String()),
			Runner:                  finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			PostDutyRunnerStateRoot: postDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     postDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:                  finishRunner(testingutils.ProposerBlindedBlockRunner(ks), testingutils.TestingProposerDutyNextEpochV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			PostDutyRunnerStateRoot: postDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*StartNewRunnerDutySpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
