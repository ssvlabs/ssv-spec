package postconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialInvalidSigQuorumThenValidQuorum tests a runner receiving a partially invalid message (due to wrong sigs) forming an invalid quorum, then receiving a valid message forming a valid quorum, terminating successfully
func PartialInvalidSigQuorumThenValidQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)
	expectedErrorCode := types.ReconstructSignatureErrorCode

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus partial invalid sig quorum then valid quorum",
		testdoc.PostConsensusPartialInvalidSignatureQuorumThenValidQuorumDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	//// Aggregator committee duty
	//multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
	//	Name: "sync committee contribution",
	//	Runner: decideAggregatorCommitteeRunner(
	//		testingutils.AggregatorCommitteeRunner(ks),
	//		testingutils.TestingSyncCommitteeContributionDuty,
	//		testingutils.TestSyncCommitteeContributionConsensusData,
	//	),
	//	Duty: testingutils.TestingSyncCommitteeContributionDuty,
	//	Messages: []*types.SignedSSVMessage{
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusPartiallyWrongSCCMsgForKeySet(ksMap, 1, spec.DataVersionPhase0, false, true))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 2, spec.DataVersionPhase0))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 3, spec.DataVersionPhase0))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 4, spec.DataVersionPhase0))),
	//	},
	//	DontStartDuty: true,
	//})
	//
	//for _, version := range testingutils.SupportedAggregatorVersions {
	//	multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
	//		{
	//			Name: fmt.Sprintf("aggregator (%s)", version.String()),
	//			Runner: decideAggregatorCommitteeRunner(
	//				testingutils.AggregatorCommitteeRunner(ks),
	//				testingutils.TestingAggregatorDuty(version),
	//				testingutils.TestAggregatorConsensusData(version),
	//			),
	//			Duty: testingutils.TestingAggregatorDuty(version),
	//			Messages: []*types.SignedSSVMessage{
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusPartiallyWrongAggMsgForKeySet(ksMap, 1, version, false, true))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 2, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 3, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 4, version))),
	//			},
	//			DontStartDuty: true,
	//		},
	//		{
	//			Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
	//			Runner: decideAggregatorCommitteeRunner(
	//				testingutils.AggregatorCommitteeRunner(ks),
	//				testingutils.TestingAggregatorCommitteeDutyMixed(version),
	//				testingutils.TestAggregatorCommitteeConsensusData(version),
	//			),
	//			Duty: testingutils.TestingAggregatorCommitteeDutyMixed(version),
	//			Messages: []*types.SignedSSVMessage{
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongAggAndSCCMsgForKeySet(ksMap, 1, version, false, true))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 2, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 3, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 4, version))),
	//			},
	//			DontStartDuty: true,
	//		},
	//	}...)
	//}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunnerWithShareMap(shareMap),
					testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongBeaconSigAttestationMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 4, version))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
				ExpectedErrorCode:      expectedErrorCode,
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunnerWithShareMap(shareMap),
					testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongBeaconSigSyncCommitteeMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 4, version))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
				ExpectedErrorCode:      expectedErrorCode,
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunnerWithShareMap(shareMap),
					testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongBeaconSigAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 4, version))),
				},
				BeaconBroadcastedRoots: append(
					testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
					testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version)...),
				DontStartDuty:     true,
				ExpectedErrorCode: expectedErrorCode,
			},
		}...)
	}

	return multiSpecTest
}
