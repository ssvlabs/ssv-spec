package committeesingleduty

import (
	"crypto/sha256"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PastMessageDutyFinished tests a valid proposal past msg for a duty that has finished
func PastMessageDutyFinished() tests.SpecTest {

	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	ks := testingutils.Testing4SharesSet()

	decidedValue := testingutils.TestBeaconVoteByts
	msgID := testingutils.CommitteeMsgID(ks)
	pastHeight := qbft.Height(10)

	attestationMessages := []*types.SignedSSVMessage{
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 1, pastHeight))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 2, pastHeight))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySet(ksMap, 3, pastHeight))),
	}

	syncCommitteeMessages := []*types.SignedSSVMessage{
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
			testingutils.PostConsensusSyncCommitteeMsgForKeySetWithSlot(ksMap, 1, phase0.Slot(pastHeight)))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
			testingutils.PostConsensusSyncCommitteeMsgForKeySetWithSlot(ksMap, 2, phase0.Slot(pastHeight)))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
			testingutils.PostConsensusSyncCommitteeMsgForKeySetWithSlot(ksMap, 3, phase0.Slot(pastHeight)))),
	}

	attestationAndSyncCommitteeMessages := []*types.SignedSSVMessage{
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
			testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 1, pastHeight))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
			testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 2, pastHeight))),
		testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
			testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySet(ksMap, 3, pastHeight))),
	}

	bumpHeight := func(c *ssv.Committee, previousDuty types.Duty, postConsensusMessages []*types.SignedSSVMessage) *ssv.Committee {

		err := c.StartDuty(previousDuty.(*types.CommitteeDuty))
		if err != nil {
			panic(err)
		}

		happyFlowMessages := []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithIdentifierAndFullData(
				ks.OperatorKeys[1], types.OperatorID(1), msgID, decidedValue,
				pastHeight),
			testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, qbft.FirstRound, pastHeight, msgID, sha256.Sum256(decidedValue)),
			testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, qbft.FirstRound, pastHeight, msgID, sha256.Sum256(decidedValue)),
			testingutils.TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, qbft.FirstRound, pastHeight, msgID, sha256.Sum256(decidedValue)),

			testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, qbft.FirstRound, pastHeight, msgID, sha256.Sum256(decidedValue)),
			testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, qbft.FirstRound, pastHeight, msgID, sha256.Sum256(decidedValue)),
			testingutils.TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, qbft.FirstRound, pastHeight, msgID, sha256.Sum256(decidedValue)),
		}

		happyFlowMessages = append(happyFlowMessages, postConsensusMessages...)

		for _, msg := range happyFlowMessages {
			err := c.ProcessMessage(msg)
			if err != nil {
				panic(err)
			}
		}

		// Erase broadcasted messages and roots due to test setup
		c.Runners[previousDuty.DutySlot()].GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs = make([]*types.SignedSSVMessage, 0)
		c.Runners[previousDuty.DutySlot()].GetBeaconNode().(*testingutils.TestingBeaconNode).BroadcastedRoots = make([]phase0.Root, 0)

		return c
	}

	pastProposalMsgF := func() *types.SignedSSVMessage {
		fullData := decidedValue
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     pastHeight,
			Round:      qbft.FirstRound,
			Identifier: msgID,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.OperatorKeys[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	expectedError := "failed processing consensus message: could not process msg: invalid signed message: proposal is not valid with current state"

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name: "past msg duty finished",
		Tests: []*committee.CommitteeSpecTest{
			{
				Name: fmt.Sprintf("%v attestation", numValidators),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap),
					testingutils.TestingCommitteeAttesterDuty(phase0.Slot(pastHeight), validatorsIndexList),
					attestationMessages),
				Input: []interface{}{
					testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name: fmt.Sprintf("%v sync committee", numValidators),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap),
					testingutils.TestingCommitteeSyncCommitteeDuty(phase0.Slot(pastHeight), validatorsIndexList),
					syncCommitteeMessages),
				Input: []interface{}{
					testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name: fmt.Sprintf("%v attestation %v sync committee", numValidators, numValidators),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap),
					testingutils.TestingCommitteeDuty(phase0.Slot(pastHeight), validatorsIndexList, validatorsIndexList),
					attestationAndSyncCommitteeMessages),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, validatorsIndexList, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		},
	}

	return multiSpecTest
}
