package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostFinish tests a msg received post runner finished
func PostFinish() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// TODO: check errors
	// nolint
	finishRunner := func(runner ssv.Runner, duty *types.Duty) ssv.Runner {
		runner.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		runner.GetBaseRunner().State.Finished = true
		return runner
	}

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus post finish",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee aggregator selection proof",
				Runner: finishRunner(
					testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], 4)),
				},
				PostDutyRunnerStateRoot: postFinishSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     postFinishSyncCommitteeContributionSC().ExpectedState,
				DontStartDuty:           true,
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           "failed processing sync committee selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name: "aggregator selection proof",
				Runner: finishRunner(
					testingutils.AggregatorRunner(ks),
					&testingutils.TestingAggregatorDuty,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAggregator(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], 4)),
				},
				PostDutyRunnerStateRoot: postFinishAggregatorSC().Root(),
				PostDutyRunnerState:     postFinishAggregatorSC().ExpectedState,
				DontStartDuty:           true,
				OutputMessages:          []*types.PartialSignatureMessages{},
				ExpectedError:           "failed processing selection proof message: invalid pre-consensus message: no running duty",
			},
			{
				Name: "validator registration",
				Runner: finishRunner(
					testingutils.ValidatorRegistrationRunner(ks),
					&testingutils.TestingValidatorRegistrationDuty,
				),
				Duty: &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgValidatorRegistration(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: postFinishValidatorRegistrationSC().Root(),
				PostDutyRunnerState:     postFinishValidatorRegistrationSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing validator registration message: invalid pre-consensus message: no running duty",
			},
			{
				Name: "voluntary exit",
				Runner: finishRunner(
					testingutils.VoluntaryExitRunner(ks),
					&testingutils.TestingVoluntaryExitDuty,
				),
				Duty: &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgVoluntaryExit(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: postFinishVoluntaryExitSC().Root(),
				PostDutyRunnerState:     postFinishVoluntaryExitSC().ExpectedState,
				OutputMessages:          []*types.PartialSignatureMessages{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing voluntary exit message: invalid pre-consensus message: no running duty",
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("randao (%s)", version.String()),
			Runner: finishRunner(
				testingutils.ProposerRunner(ks),
				testingutils.TestingProposerDutyV(version),
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, version)),
			},
			PostDutyRunnerStateRoot: postFinishProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishProposerSC(version).ExpectedState,
			DontStartDuty:           true,
			OutputMessages:          []*types.PartialSignatureMessages{},
			ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name: fmt.Sprintf("randao blinded block (%s)", version.String()),
			Runner: finishRunner(
				testingutils.ProposerBlindedBlockRunner(ks),
				testingutils.TestingProposerDutyV(version),
			),
			Duty: testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, version)),
			},
			PostDutyRunnerStateRoot: postFinishBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     postFinishBlindedProposerSC(version).ExpectedState,
			DontStartDuty:           true,
			OutputMessages:          []*types.PartialSignatureMessages{},
			ExpectedError:           "failed processing randao message: invalid pre-consensus message: no running duty",
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
