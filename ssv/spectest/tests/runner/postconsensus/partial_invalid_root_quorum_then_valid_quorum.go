package postconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialInvalidRootQuorumThenValidQuorum tests a runner receiving a partially invalid message (due to wrong roots) forming an invalid quorum, then receiving a valid message forming a valid quorum, terminating successfully
func PartialInvalidRootQuorumThenValidQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)
	valIdxs := make([]int, 0)
	for valIdx := range ksMap {
		valIdxs = append(valIdxs, int(valIdx))
	}

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus partial invalid root quorum then valid quorum",
		testdoc.PostConsensusPartialInvalidRootQuorumThenValidQuorumDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	//// Aggregator committee duty
	//sccDuty := testingutils.TestingSyncCommitteeContributorDutyForValidators(spec.DataVersionPhase0, valIdxs)
	//multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
	//	Name: "sync committee contribution",
	//	Runner: decideAggregatorCommitteeRunner(
	//		testingutils.AggregatorCommitteeRunner(ks),
	//		sccDuty,
	//		testingutils.TestSyncCommitteeContributionConsensusDataForDuty(sccDuty),
	//	),
	//	Duty: sccDuty,
	//	Messages: []*types.SignedSSVMessage{
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusPartiallyWrongSCCMsgForKeySet(ksMap, 1, spec.DataVersionPhase0, true, false))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 2, spec.DataVersionPhase0))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 3, spec.DataVersionPhase0))),
	//		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 4, spec.DataVersionPhase0))),
	//	},
	//	BeaconBroadcastedRoots: syncCommitteeRootsForKSMap(ksMap),
	//	DontStartDuty:          true,
	//	ExpectedErrorCode:      types.PostConsensusQuorumWithInvalidSignatures,
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
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusPartiallyWrongAggMsgForKeySet(ksMap, 1, version, true, false))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 2, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 3, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 4, version))),
	//			},
	//			BeaconBroadcastedRoots: aggregatorRootForKSMap(ksMap, version),
	//			DontStartDuty:          true,
	//			ExpectedErrorCode:      types.PostConsensusQuorumWithInvalidSignatures,
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
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongAggAndSCCMsgForKeySet(ksMap, 1, version, true, false))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 2, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 3, version))),
	//				testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 4, version))),
	//			},
	//			BeaconBroadcastedRoots: aggregatorAndSyncCommitteeRootsForKSMap(ksMap, version),
	//			DontStartDuty:          true,
	//			ExpectedErrorCode:      types.PostConsensusQuorumWithInvalidSignatures,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongRootAttestationMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 4, version))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongRootSyncCommitteeMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 4, version))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongRootAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 4, version))),
				},
				BeaconBroadcastedRoots: append(
					testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
					testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version)...),
				DontStartDuty: true,
			},
		}...)
	}

	return multiSpecTest
}
