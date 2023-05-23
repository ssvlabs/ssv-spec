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

// NotDecided tests starting duty before finished or decided
func NotDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check error
	// nolint
	startRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.FirstHeight
		return r
	}

	expectedErr := "previous instance hasn't Decided"
	multiSpecTest := &MultiStartNewRunnerDutySpecTest{
		Name: "new duty not decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  startRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionNexEpochDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "bb6cec38c88225ab7f6f5020165695964f90b11e131421d37632caf34820a502",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:                    "sync committee",
				Runner:                  startRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "fe3c3934ed1a0cc6d3dffa04cbe840abf363e64fc4a15e1d7cdf947fe83adb70",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedErr,
			},
			{
				Name:                    "aggregator",
				Runner:                  startRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDutyNextEpoch),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "28ae79d6c896dda7d2511617bbc606a7880d0d7a029840630589296b82ccb258",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:                    "proposer",
				Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix)),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionBellatrix),
				PostDutyRunnerStateRoot: "25ea4647a66e49db7f8d0fccb7ba07cdeb7d10ff129dbb68b9d499093cde8d6a",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:                    "attester",
				Runner:                  startRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "420e7c638b98b00c0af8e2ac641f670d820d3fa4b583eb678cb7de36688cb087",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedErr,
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer (%s)", version.String()),
			Runner:                  startRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyNextEpochV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			PostDutyRunnerStateRoot: notDecidedProposerSC(version).Root(),
			PostDutyRunnerState:     notDecidedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
			ExpectedError: expectedErr,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *StartNewRunnerDutySpecTest {
		return &StartNewRunnerDutySpecTest{
			Name:                    fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner:                  startRunner(testingutils.ProposerBlindedBlockRunner(ks), testingutils.TestingProposerDutyNextEpochV(version)),
			Duty:                    testingutils.TestingProposerDutyNextEpochV(version),
			PostDutyRunnerStateRoot: notDecidedBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     notDecidedBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.SignedPartialSignatureMessage{
				testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
			ExpectedError: expectedErr,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*StartNewRunnerDutySpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
