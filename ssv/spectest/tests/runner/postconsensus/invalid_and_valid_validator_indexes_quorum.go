package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidAndValidValidatorIndexesQuorum tests a quorum of post consensus messages with both an invalid and a valid validator index
func InvalidAndValidValidatorIndexesQuorum() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	validatorsIndex := []phase0.ValidatorIndex{testingutils.TestingWrongValidatorIndex, testingutils.TestingValidatorIndex}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus invalid and valid validator index quorum",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "attester",
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterDuty,
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForValidatorsIndex(ks.Shares[1], 1, testingutils.TestingDutySlot, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForValidatorsIndex(ks.Shares[2], 2, testingutils.TestingDutySlot, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForValidatorsIndex(ks.Shares[3], 3, testingutils.TestingDutySlot, validatorsIndex))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestationForValidatorIndex(ks, testingutils.TestingValidatorIndex)),
				},
				DontStartDuty: true,
			},
			{
				Name: "sync committee",
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingSyncCommitteeDuty,
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForValidatorsIndex(ks.Shares[1], 1, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForValidatorsIndex(ks.Shares[2], 2, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusSyncCommitteeMsgForValidatorsIndex(ks.Shares[3], 3, validatorsIndex))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRootForValidatorIndex(ks, testingutils.TestingValidatorIndex)),
				},
				DontStartDuty: true,
			},
			{
				Name: "attester and sync committee",
				Runner: decideCommitteeRunner(
					testingutils.CommitteeRunner(ks),
					testingutils.TestingAttesterAndSyncCommitteeDuties,
					&testingutils.TestBeaconVote,
				),
				Duty: testingutils.TestingAttesterAndSyncCommitteeDuties,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex(ks.Shares[1], 1, testingutils.TestingDutySlot, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex(ks.Shares[2], 2, testingutils.TestingDutySlot, validatorsIndex))),
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationAndSyncCommitteeMsgForValidatorsIndex(ks.Shares[3], 3, testingutils.TestingDutySlot, validatorsIndex))),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: []string{
					testingutils.GetSSZRootNoError(testingutils.TestingSignedAttestationForValidatorIndex(ks, testingutils.TestingValidatorIndex)),
					testingutils.GetSSZRootNoError(testingutils.TestingSignedSyncCommitteeBlockRootForValidatorIndex(ks, testingutils.TestingValidatorIndex)),
				},
				DontStartDuty: true,
			},
		},
	}
}
