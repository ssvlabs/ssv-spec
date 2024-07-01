package postconsensus

import (
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

	attestationCommitteeDuty := testingutils.TestingCommitteeDutyWithMixedCommitteeIndexes(testingutils.TestingDutySlot, validatorsIndexList, nil)
	syncCommitteeCommitteeDuty := testingutils.TestingCommitteeDutyWithMixedCommitteeIndexes(testingutils.TestingDutySlot, nil, validatorsIndexList)
	attestationAndSyncCommitteeCommitteeDuty := testingutils.TestingCommitteeDutyWithMixedCommitteeIndexes(testingutils.TestingDutySlot, validatorsIndexList, validatorsIndexList)

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "mixed committees",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "attester",
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
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(attestationCommitteeDuty, ksMap),
				DontStartDuty:          true,
			},
			{
				Name: "sync committee",
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
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(syncCommitteeCommitteeDuty, ksMap),
				DontStartDuty:          true,
			},
			{
				Name: "attester and sync committee",
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
				OutputMessages:         []*types.PartialSignatureMessages{},
				BeaconBroadcastedRoots: testingutils.TestingSignedCommitteeBeaconObjectSSZRoot(attestationAndSyncCommitteeCommitteeDuty, ksMap),
				DontStartDuty:          true,
			},
		},
	}

	return multiSpecTest
}
