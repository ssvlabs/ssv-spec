package runner

import (
	"encoding/hex"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	comparable2 "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
	ssz "github.com/ferranbt/fastssz"
)

func getSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

func hexDecodeNoErr(h string) []byte {
	ret, err := hex.DecodeString(h)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

func runnerHappyFlowStates(ks *testingutils.TestKeySet) []ssv.Runner {
	syncCommitteeContrib := testingutils.SyncCommitteeContributionRunner(ks)
	syncCommitteeContrib.GetBaseRunner().State = &ssv.State{
		Finished:               true,
		DecidedValue:           comparable2.FixIssue178(testingutils.TestSyncCommitteeContributionConsensusData, spec.DataVersionPhase0),
		StartingDuty:           &testingutils.TestSyncCommitteeContributionConsensusData.Duty,
		PreConsensusContainer:  ssv.NewPartialSigContainer(3),
		PostConsensusContainer: ssv.NewPartialSigContainer(3),
	}
	ssvcomparable.SetMessagesInContainer(
		syncCommitteeContrib.GetBaseRunner().State.PreConsensusContainer,
		testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution)[:3])
	ssvcomparable.SetMessagesInContainer(
		syncCommitteeContrib.GetBaseRunner().State.PostConsensusContainer,
		[]*types.SSVMessage{
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
		})
	syncCommitteeContrib.GetBaseRunner().State.RunningInstance = &qbft.Instance{
		StartValue: comparable2.NoErrorEncoding(comparable2.FixIssue178(testingutils.TestSyncCommitteeContributionConsensusData, spec.DataVersionBellatrix)),
		State: &qbft.State{
			Share:  testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:     syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
			Round:  qbft.FirstRound,
			Height: qbft.FirstHeight,
			ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
				ks.Shares[1], types.OperatorID(1),
				syncCommitteeContrib.GetBaseRunner().QBFTController.Identifier,
				testingutils.TestSyncCommitteeContributionConsensusDataByts,
			),
			LastPreparedRound: 1,
			LastPreparedValue: testingutils.TestSyncCommitteeContributionConsensusDataByts,
			Decided:           true,
			DecidedValue:      testingutils.TestSyncCommitteeContributionConsensusDataByts,
		},
	}
	qbftcomparable.SetMessages(
		syncCommitteeContrib.GetBaseRunner().State.RunningInstance,
		testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution)[3:10],
	)
	syncCommitteeContrib.GetBaseRunner().QBFTController.StoredInstances = append(syncCommitteeContrib.GetBaseRunner().QBFTController.StoredInstances, syncCommitteeContrib.GetBaseRunner().State.RunningInstance)

	syncCommittee := testingutils.SyncCommitteeRunner(ks)

	aggregator := testingutils.AggregatorRunner(ks)

	proposer := testingutils.ProposerRunner(ks)

	blindedProposer := testingutils.ProposerBlindedBlockRunner(ks)

	attester := testingutils.AttesterRunner(ks)

	validatorRegistration := testingutils.ValidatorRegistrationRunner(ks)

	return []ssv.Runner{
		syncCommitteeContrib,
		syncCommittee,
		aggregator,
		proposer,
		blindedProposer,
		attester,
		validatorRegistration,
	}
}

// FullHappyFlow  tests a full runner happy flow
func FullHappyFlow() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()

	// register runners
	runnerStates := runnerHappyFlowStates(ks)
	roots := make([]string, 0)
	for _, runner := range runnerStates {
		r, err := runner.GetRoot()
		if err != nil {
			panic(err.Error())
		}
		roots = append(roots, hex.EncodeToString(r[:]))
		ssvcomparable.RootRegister[hex.EncodeToString(r[:])] = runner
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "full happy flow",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: append(
					[]*types.SSVMessage{ // pre consensus
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
						testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					},
					append(
						// consensus
						testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution),
						[]*types.SSVMessage{ // post consensus
							testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks)),
							testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[2], 2, ks)),
							testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[3], 3, ks)),
						}...,
					)...,
				),
				PostDutyRunnerStateRoot: roots[0], //"4987127ad389bb9d21500d447686f135a19f59ae10192e82bf052278853ad3d1",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[0], testingutils.TestingContributionProofsSigned[0], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[1], testingutils.TestingContributionProofsSigned[1], ks)),
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeContributions(testingutils.TestingSyncCommitteeContributions[2], testingutils.TestingContributionProofsSigned[2], ks)),
				},
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "48c73f57659b69131467ef133ccb35d7de2fe96438d30bfa2b5ea63b19ead011",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRoot(ks)),
				},
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: append(
					[]*types.SSVMessage{ // pre consensus
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
						testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
					},
					append(
						testingutils.SSVDecidingMsgs(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator), // consensus
						[]*types.SSVMessage{ // post consensus
							testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1)),
							testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[2], 2)),
							testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[3], 3)),
						}...,
					)...,
				),
				PostDutyRunnerStateRoot: "298bbb63d87a36eef30926c2c21baad6990db0f8fa03a83ca56b2463c7f0065c",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAggregateAndProof(ks)),
				},
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerConsensusData, ks, types.BNRoleProposer), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "76812c0f14ff09067547e9528730749b0c0090d1a4872689a0b8480d7b538884",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "proposer blinded block",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   &testingutils.TestingProposerDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestProposerBlindedBlockConsensusData, ks, types.BNRoleProposer), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[1], 1)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[2], 2)),
						testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsg(ks.Shares[3], 3)),
					}...,
				),
				PostDutyRunnerStateRoot: "90755cc41b814519fd9fdd14bc82d239997ba51340c297f25f5f1552f27f66c7",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1),
					testingutils.PostConsensusProposerMsg(ks.Shares[1], 1),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedBeaconBlock(ks)),
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: append(
					testingutils.SSVDecidingMsgs(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester), // consensus
					[]*types.SSVMessage{ // post consensus
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
						testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
					}...,
				),
				PostDutyRunnerStateRoot: "9d55ff5721b21c5b99dd4b4bacb0acda0b674112fe3cec55cc6aeb04ad5dc2fc",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
				BeaconBroadcastedRoots: []string{
					getSSZRootNoError(testingutils.TestingSignedAttestation(ks)),
				},
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[2], 2)),
					testingutils.SSVMsgValidatorRegistration(nil, testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[3], 3)),
				},
				PostDutyRunnerStateRoot: "f36c8b537afaba0894dbc8c87cb94466d8ac2623e9283f1c584e3d544b5f2b88",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
		},
	}
}
