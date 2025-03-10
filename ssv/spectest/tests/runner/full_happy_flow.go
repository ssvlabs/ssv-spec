package runner

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FullHappyFlow tests a full runner happy flow
func FullHappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "full happy flow",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					// consensus
					testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.RoleSyncCommitteeContribution),
					[]*types.SignedSSVMessage{ // post consensus
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks))),
					}...,
				),
				PostDutyRunnerStateRoot: fullHappyFlowSyncCommitteeContributionSC().Root(), //"4987127ad389bb9d21500d447686f135a19f59ae10192e82bf052278853ad3d1",
				PostDutyRunnerState:     fullHappyFlowSyncCommitteeContributionSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3))),
				},
				PostDutyRunnerStateRoot: fullHappyFlowValidatorRegistrationSC().Root(), // "ec573732e70b70808972c43acb5ead6443cff06ba30d8abb51e37ac82ffe0727",
				PostDutyRunnerState:     fullHappyFlowValidatorRegistrationSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedValidatorRegistration(ks)),
				},
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgVoluntaryExit(nil, testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[3], 3))),
				},
				PostDutyRunnerStateRoot: fullHappyFlowVoluntaryExitSC().Root(), // "ec573732e70b70808972c43acb5ead6443cff06ba30d8abb51e37ac82ffe0727",
				PostDutyRunnerState:     fullHappyFlowVoluntaryExitSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedVoluntaryExit(ks)),
				},
			},
			{
				Name:   "preconfirmation",
				Runner: testingutils.PreconfRunner(ks),
				Duty:   &testingutils.TestingPreconfDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPreconf(nil, testingutils.PreConsensusPreconfMsg(ks.Shares[1], 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPreconf(nil, testingutils.PreConsensusPreconfMsg(ks.Shares[2], 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgPreconf(nil, testingutils.PreConsensusPreconfMsg(ks.Shares[3], 3))),
				},
				PostDutyRunnerStateRoot: fullHappyFlowPreconfSC().Root(),
				PostDutyRunnerState:     fullHappyFlowPreconfSC().ExpectedState,
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PreConsensusPreconfMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				// preconfirmation duty does not broadcast to beacon
			},
		},
	}

	for _, version := range testingutils.SupportedAggregatorVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("aggregator (%s)", version.String()),
			Runner: testingutils.AggregatorRunner(ks),
			Duty:   testingutils.TestingAggregatorDuty(version),
			Messages: append(
				[]*types.SignedSSVMessage{ // pre consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3, version))),
				},
				append(
					testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData(version), ks, types.RoleAggregator), // consensus
					[]*types.SignedSSVMessage{ // post consensus
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3, version))),
					}...,
				)...,
			),
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1, version),
				testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks, version)),
			},
		})
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		slot := testingutils.TestingDutySlotV(version)
		height := qbft.Height(slot)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name:   fmt.Sprintf("attestation (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty(version),
				Messages: append(
					// consensus
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
					[]*types.SignedSSVMessage{ // post consensus
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, version))),
					}...,
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, version),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObject(ks, version)),
				},
			},
			{
				Name:   fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty(version),
				Messages: append(
					// consensus
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
					[]*types.SignedSSVMessage{ // post consensus
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3, version))),
					}...,
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1, version),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks, version)),
				},
			},
			{
				Name:   fmt.Sprintf("attestation and sync committee (%s)", version.String()),
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: append(
					// consensus
					testingutils.SSVDecidingMsgsForCommitteeRunner(&testingutils.TestBeaconVote, ks, height),
					[]*types.SignedSSVMessage{ // post consensus
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[2], 2, version))),
						testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[3], 3, version))),
					}...,
				),
				OutputMessages: []*types.PartialSignatureMessages{
					testingutils.PostConsensusAttestationAndSyncCommitteeMsg(ks.Shares[1], 1, version),
				},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObject(ks, version)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks, version)),
				},
			},
		}...)
	}

	// proposerV creates a test specification for versioned proposer.
	proposerV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer (%s)", version.String()),
			Runner: testingutils.ProposerRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(version), ks, types.RoleProposer), // consensus
				[]*types.SignedSSVMessage{ // post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				}...,
			),
			PostDutyRunnerStateRoot: fullHappyFlowProposerSC(version).Root(),
			PostDutyRunnerState:     fullHappyFlowProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version)),
			},
		}
	}

	// proposerBlindedV creates a test specification for versioned proposer with blinded block.
	proposerBlindedV := func(version spec.DataVersion) *tests.MsgProcessingSpecTest {
		return &tests.MsgProcessingSpecTest{
			Name:   fmt.Sprintf("proposer blinded block (%s)", version.String()),
			Runner: testingutils.ProposerBlindedBlockRunner(ks),
			Duty:   testingutils.TestingProposerDutyV(version),
			Messages: append(
				testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(version), ks, types.RoleProposer), // consensus
				[]*types.SignedSSVMessage{ // post consensus
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))),
				}...,
			),
			PostDutyRunnerStateRoot: fullHappyFlowBlindedProposerSC(version).Root(),
			PostDutyRunnerState:     fullHappyFlowBlindedProposerSC(version).ExpectedState,
			OutputMessages: []*types.PartialSignatureMessages{
				testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, version),
				testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version),
			},
			BeaconBroadcastedRoots: []string{
				testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, version)),
			},
		}
	}

	for _, v := range testingutils.SupportedBlockVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{proposerV(v), proposerBlindedV(v)}...)
	}

	return multiSpecTest
}
