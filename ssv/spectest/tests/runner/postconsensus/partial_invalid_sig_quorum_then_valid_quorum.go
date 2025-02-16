package postconsensus

import (
	"fmt"

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
	expectedError := "got post-consensus quorum but it has invalid signatures: could not reconstruct beacon sig: failed to verify reconstruct signature: could not reconstruct a valid signature"

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name:  "post consensus partial invalid sig quorum then valid quorum",
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
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusPartiallyWrongBeaconSigAttestationMsgForKeySet(ksMap, 1, version))),

					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 2, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 3, version))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 4, version))),
				},
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
				ExpectedError:          expectedError,
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
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version),
				DontStartDuty:          true,
				ExpectedError:          expectedError,
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
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: append(
					testingutils.TestingSignedAttestationResponseSSZRootForKeyMap(ksMap, version),
					testingutils.TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap(ksMap, version)...),
				DontStartDuty: true,
				ExpectedError: expectedError,
			},
		}...)
	}

	return multiSpecTest
}
