package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MixedCommittees tests a committee runner with duties with different CommitteeIndex
func MixedCommittees() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	multiSpecTest := tests.NewMultiMsgProcessingSpecTest(
		"mixed committees",
		testdoc.PostConsensusMixedCommitteesDoc,
		[]*tests.MsgProcessingSpecTest{},
		ks,
	)

	for _, version := range testingutils.SupportedAttestationVersions {

		attestationCommitteeDuty := testingutils.TestingCommitteeDutyWithMixedCommitteeIndexes(validatorsIndexList, nil, version)
		syncCommitteeCommitteeDuty := testingutils.TestingCommitteeDutyWithMixedCommitteeIndexes(nil, validatorsIndexList, version)
		attestationAndSyncCommitteeCommitteeDuty := testingutils.TestingCommitteeDutyWithMixedCommitteeIndexes(validatorsIndexList, validatorsIndexList, version)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunnerWithShareMap(shareMap),
					attestationCommitteeDuty,
					&testingutils.TestBeaconVote,
				),
				Duty: attestationCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(attestationCommitteeDuty, ksMap, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(attestationCommitteeDuty, ksMap, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(attestationCommitteeDuty, ksMap, 3))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(attestationCommitteeDuty, ksMap, version),
				DontStartDuty:          true,
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunnerWithShareMap(shareMap),
					syncCommitteeCommitteeDuty,
					&testingutils.TestBeaconVote,
				),
				Duty: syncCommitteeCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(syncCommitteeCommitteeDuty, ksMap, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(syncCommitteeCommitteeDuty, ksMap, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(syncCommitteeCommitteeDuty, ksMap, 3))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(syncCommitteeCommitteeDuty, ksMap, version),
				DontStartDuty:          true,
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunnerWithShareMap(shareMap),
					attestationAndSyncCommitteeCommitteeDuty,
					&testingutils.TestBeaconVote,
				),
				Duty: attestationAndSyncCommitteeCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(attestationAndSyncCommitteeCommitteeDuty, ksMap, 1))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(attestationAndSyncCommitteeCommitteeDuty, ksMap, 2))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusCommitteeMsgForDuty(attestationAndSyncCommitteeCommitteeDuty, ksMap, 3))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(attestationAndSyncCommitteeCommitteeDuty, ksMap, version),
				DontStartDuty:          true,
			},
		}...)
	}

	// Aggregator committee duty
	sccCommitteeDuty := testingutils.TestingAggCommitteeDutyWithMixedCommitteeIndexes(validatorsIndexList, nil, spec.DataVersionPhase0)
	multiSpecTest.Tests = append(multiSpecTest.Tests, &tests.MsgProcessingSpecTest{
		Name: "sync committee contributor",
		Runner: decideAggregatorCommitteeRunner(
			testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
			sccCommitteeDuty,
			testingutils.TestAggregatorCommitteeConsensusDataForDuty(sccCommitteeDuty, spec.DataVersionPhase0),
		),
		Duty: sccCommitteeDuty,
		Messages: []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(sccCommitteeDuty, ksMap, 1, spec.DataVersionPhase0))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(sccCommitteeDuty, ksMap, 2, spec.DataVersionPhase0))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(sccCommitteeDuty, ksMap, 3, spec.DataVersionPhase0))),
		},
		BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(sccCommitteeDuty, ksMap, spec.DataVersionPhase0),
		DontStartDuty:          true,
	})

	for _, version := range testingutils.SupportedAggregatorVersions {
		aggCommitteeDuty := testingutils.TestingAggCommitteeDutyWithMixedCommitteeIndexes(nil, validatorsIndexList, version)
		mixedAggCommitteeDuty := testingutils.TestingAggCommitteeDutyWithMixedCommitteeIndexes(validatorsIndexList, validatorsIndexList, version)
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{
			{
				Name: fmt.Sprintf("aggregator (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
					aggCommitteeDuty,
					testingutils.TestAggregatorCommitteeConsensusDataForDuty(aggCommitteeDuty, version),
				),
				Duty: aggCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(aggCommitteeDuty, ksMap, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(aggCommitteeDuty, ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(aggCommitteeDuty, ksMap, 3, version))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(aggCommitteeDuty, ksMap, version),
				DontStartDuty:          true,
			},
			{
				Name: fmt.Sprintf("aggregator committee mixed (%s)", version.String()),
				Runner: decideAggregatorCommitteeRunner(
					testingutils.AggregatorCommitteeRunnerWithShareMap(shareMap),
					mixedAggCommitteeDuty,
					testingutils.TestAggregatorCommitteeConsensusDataForDuty(mixedAggCommitteeDuty, version),
				),
				Duty: mixedAggCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(mixedAggCommitteeDuty, ksMap, 1, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(mixedAggCommitteeDuty, ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgAggregatorCommittee(ks, nil, testingutils.PostConsensusAggregatorCommitteeMsgForDuty(mixedAggCommitteeDuty, ksMap, 3, version))),
				},
				BeaconBroadcastedRoots: testingutils.TestingSignedAggregatorCommitteeBeaconObjectSSZRoot(mixedAggCommitteeDuty, ksMap, version),
				DontStartDuty:          true,
			},
		}...)
	}
	return multiSpecTest
}
