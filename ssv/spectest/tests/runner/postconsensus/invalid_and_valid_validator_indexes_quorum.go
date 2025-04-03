package postconsensus

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidAndValidValidatorIndexesQuorum tests a quorum of post consensus messages with both an invalid and a valid validator index
func InvalidAndValidValidatorIndexesQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	validatorsIndex := []phase0.ValidatorIndex{testingutils.TestingWrongValidatorIndex, testingutils.TestingValidatorIndex}

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name:  "post consensus invalid and valid validator index quorum",
		Tests: []*tests.MsgProcessingSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*tests.MsgProcessingSpecTest{

			{
				Name: fmt.Sprintf("attester (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForValidatorsIndex(ks.Shares[1], 1, version, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForValidatorsIndex(ks.Shares[2], 2, version, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForValidatorsIndex(ks.Shares[3], 3, version, validatorsIndex))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObjectForValidatorIndex(ks, version, testingutils.TestingValidatorIndex)),
				},
				DontStartDuty: true,
			},
			{
				Name: fmt.Sprintf("sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForValidatorsIndex(ks.Shares[1], 1, version, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForValidatorsIndex(ks.Shares[2], 2, version, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForValidatorsIndex(ks.Shares[3], 3, version, validatorsIndex))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRootForValidatorIndex(ks, testingutils.TestingValidatorIndex, version)),
				},
				DontStartDuty: true,
			},
			{
				Name: fmt.Sprintf("attester and sync committee (%s)", version.String()),
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties(version),
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties(version),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex(ks.Shares[1], 1, version, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex(ks.Shares[2], 2, version, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex(ks.Shares[3], 3, version, validatorsIndex))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingAttestationResponseBeaconObjectForValidatorIndex(ks, version, testingutils.TestingValidatorIndex)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRootForValidatorIndex(ks, testingutils.TestingValidatorIndex, version)),
				},
				DontStartDuty: true,
			},
		}...)
	}

	return multiSpecTest
}
