package preconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Quorum7Operators  tests a quorum of valid PartialSignatureMessages (7 operators)
func Quorum7Operators() tests.SpecTest {
	ks := testingutils.Testing7SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus quorum 7 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgSyncCommitteeContribution(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgSyncCommitteeContribution(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgSyncCommitteeContribution(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgSyncCommitteeContribution(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: quorum7OperatorsSyncCommitteeContributionSC().Root(),
				PostDutyRunnerState:     quorum7OperatorsSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgAggregator(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgAggregator(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgAggregator(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgAggregator(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgAggregator(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: quorum7OperatorsAggregatorSC().Root(),
				PostDutyRunnerState:     quorum7OperatorsAggregatorSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgValidatorRegistration(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgValidatorRegistration(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgValidatorRegistration(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgValidatorRegistration(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgValidatorRegistration(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: quorum7OperatorsValidatorRegistrationSC().Root(),
				PostDutyRunnerState:     quorum7OperatorsValidatorRegistrationSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingValidatorRegistration),
				},
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SSVMsgVoluntaryExit(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgVoluntaryExit(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgVoluntaryExit(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3)),
					testingutils.SSVMsgVoluntaryExit(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[4], 4)),
					testingutils.SSVMsgVoluntaryExit(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[5], 5)),
				},
				PostDutyRunnerStateRoot: quorum7OperatorsVoluntaryExitSC().Root(),
				PostDutyRunnerState:     quorum7OperatorsVoluntaryExitSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
			},
		},
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("randao (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, version)),
				testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, version)),
				testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, version)),
			},
			PostDutyRunnerStateRoot: quorum7OperatorsProposerSC(version).Root(),
			PostDutyRunnerState:     quorum7OperatorsProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("randao blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: []*types.SignedSSVMessage{
				testingutils.SSVMsgProposer(1, ks.NetworkKeys[1], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], 1, version)),
				testingutils.SSVMsgProposer(2, ks.NetworkKeys[2], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], 2, version)),
				testingutils.SSVMsgProposer(3, ks.NetworkKeys[3], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], 3, version)),
				testingutils.SSVMsgProposer(4, ks.NetworkKeys[4], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[4], 4, version)),
				testingutils.SSVMsgProposer(5, ks.NetworkKeys[5], nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[5], 5, version)),
			},
			PostDutyRunnerStateRoot: quorum7OperatorsBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     quorum7OperatorsBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version), // broadcasts when starting a new duty
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
