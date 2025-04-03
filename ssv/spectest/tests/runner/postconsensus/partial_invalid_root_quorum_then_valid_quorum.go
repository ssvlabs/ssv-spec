package postconsensus

import (
	"fmt"

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

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name:  "post consensus partial invalid root quorum then valid quorum",
		Tests: []*tests.MsgProcessingSpecTest{},
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongRootAttestationMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 4, version))),
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
			},
			// {
			// 	Name: fmt.Sprintf("sync committee (%s)", version.String()),
			// 	Runner: decideCommitteeRunner(
			// 		testingutils.CommitteeRunnerWithShareMap(shareMap),
			// 		testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
			// 		&testingutils.TestBeaconVote,
			// 	),
			// 	Duty: testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
			// 	Messages: []*types.SignedSSVMessage{
			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongRootSyncCommitteeMsgForKeySet(ksMap, 1, version))),

			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 2, version))),
			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 3, version))),
			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForKeySet(ksMap, 4, version))),
			// 	},
			// 	OutputMessages:         []*types.PartialSignatureMessages{},
			// 	BeaconBroadcastedRoots: testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version),
			// 	DontStartDuty:          true,
			// },
			// {
			// 	Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
			// 	Runner: decideCommitteeRunner(
			// 		testingutils.CommitteeRunnerWithShareMap(shareMap),
			// 		testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
			// 		&testingutils.TestBeaconVote,
			// 	),
			// 	Duty: testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
			// 	Messages: []*types.SignedSSVMessage{
			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongRootAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, version))),

			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 2, version))),
			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 3, version))),
			// 		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 4, version))),
			// 	},
			// 	OutputMessages: []*types.PartialSignatureMessages{},
			// 	BeaconBroadcastedRoots: append(
			// 		testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
			// 		testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version)...),
			// 	DontStartDuty: true,
			// },
		}...)
	}

	return multiSpecTest
}
