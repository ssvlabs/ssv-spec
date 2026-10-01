package postconsensus

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
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
	aggCmmErrorCode := types.PostConsensusQuorumWithInvalidSignatures

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"post consensus partial invalid sig quorum then valid quorum",
		testdoc.PostConsensusPartialInvalidSignatureQuorumThenValidQuorumDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	// Aggregator committee duty
	sccDuty := testingutils.TestingAggregatorCommitteeDuty(nil, validatorsIndexList, gloas.DataVersionGloas)
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name: "sync committee contribution",
		Runner: decideAggregatorCommitteeRunner(
			testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
			sccDuty,
			testingutils.TestSyncCommitteeContributionConsensusDataForDuty(sccDuty),
		),
		Duty: sccDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusPartiallyWrongSCCMsgForKeySet(ksMap, 1, gloas.DataVersionGloas, false, true))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 2, gloas.DataVersionGloas))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 3, gloas.DataVersionGloas))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSCCMsgForKeySet(ksMap, 4, gloas.DataVersionGloas))),
		},
		DontStartDuty:          true,
		BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(sccDuty, ksMap, gloas.DataVersionGloas),
		ExpectedErrorCode:      aggCmmErrorCode,
	})

	for _, version := range testingutils.SupportedAggregatorVersions {
		aggDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, nil, version)
		aggSccDuty := testingutils.TestingAggregatorCommitteeDuty(validatorsIndexList, validatorsIndexList, version)
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name: fmt.Sprintf("aggregator (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
					aggDuty,
					testingutils.TestAggregatorCommitteeConsensusDataForDuty(aggDuty, version, nil),
				),
				Duty: aggDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusPartiallyWrongAggMsgForKeySet(ksMap, 1, version, false, true))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggMsgForKeySet(ksMap, 4, version))),
				},
				DontStartDuty:          true,
				BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(aggDuty, ksMap, version),
				ExpectedErrorCode:      aggCmmErrorCode,
			},
			{
				Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
					aggSccDuty,
					testingutils.TestAggregatorCommitteeConsensusDataForDuty(aggSccDuty, version, nil),
				),
				Duty: aggSccDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongAggAndSCCMsgForKeySet(ksMap, 1, version, false, true))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggAndSCCMsgForKeySet(ksMap, 4, version))),
				},
				DontStartDuty:          true,
				BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(aggSccDuty, ksMap, version),
				ExpectedErrorCode:      aggCmmErrorCode,
			},
		}...)
	}

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
