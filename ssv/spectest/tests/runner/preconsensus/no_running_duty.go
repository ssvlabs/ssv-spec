package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoRunningDuty tests a valid partial pre consensus msg before duty starts
func NoRunningDuty() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus no running duty",
		testdoc.PreConsensusNoRunningDutyDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
				},
				PostDutyRunnerStateRoot: noRunningDutySyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     noRunningDutySyncCommitteeContributionSC().ExpectedState,
				DontStartDuty:           true,
				ExpectedErrorCode:       types.NoRunningDutyErrorCode,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: noRunningDutyValidatorRegistrationSC().Root(),
				PostDutyRunnerState:     noRunningDutyValidatorRegistrationSC().ExpectedState,
				DontStartDuty:           true,
				ExpectedErrorCode:       types.NoRunningDutyErrorCode,
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: noRunningDutyVoluntaryExitSC().Root(),
				PostDutyRunnerState:     noRunningDutyVoluntaryExitSC().ExpectedState,
				DontStartDuty:           true,
				ExpectedErrorCode:       types.NoRunningDutyErrorCode,
			},
		},
		ks,
	)

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: testingutils.AggregatorRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
			},
			DontStartDuty:     true,
			ExpectedErrorCode: types.NoRunningDutyErrorCode,
		})
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version))),
			},
			PostDutyRunnerStateRoot: noRunningDutyProposerSC(version).Root(),
			PostDutyRunnerState:     noRunningDutyProposerSC(version).ExpectedState,
			DontStartDuty:           true,
			ExpectedErrorCode:       types.NoRunningDutyErrorCode,
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, version))),
			},
			PostDutyRunnerStateRoot: noRunningDutyBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     noRunningDutyBlindedProposerSC(version).ExpectedState,
			DontStartDuty:           true,
			ExpectedErrorCode:       types.NoRunningDutyErrorCode,
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
