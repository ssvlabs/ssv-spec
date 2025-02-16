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

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "past msg duty finished",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		pastSlot := testingutils.TestingDutySlotV(version) - 2
		pastHeight := qbft.Height(pastSlot)

		attestationMessages := []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySetWithSlot(ksMap, 1, pastSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySetWithSlot(ksMap, 2, pastSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil, testingutils.PostConsensusAttestationMsgForKeySetWithSlot(ksMap, 3, pastSlot))),
		}

		syncCommitteeMessages := []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
				testingutils.PostConsensusSyncCommitteeMsgForKeySetWithSlot(ksMap, 1, pastSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
				testingutils.PostConsensusSyncCommitteeMsgForKeySetWithSlot(ksMap, 2, pastSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
				testingutils.PostConsensusSyncCommitteeMsgForKeySetWithSlot(ksMap, 3, pastSlot))),
		}

		attestationAndSyncCommitteeMessages := []*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
				testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySetWithSlot(ksMap, 1, pastSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
				testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySetWithSlot(ksMap, 2, pastSlot))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgCommittee(ks, nil,
				testingutils.PostConsensusAttestationAndSyncCommitteeMsgForKeySetWithSlot(ksMap, 3, pastSlot))),
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

		expectedError := "failed processing consensus message: not processing consensus message since instance is already decided"

		attesterDuty := testingutils.TestingCommitteeDutyForSlot(phase0.Slot(pastHeight), validatorsIndexList, nil)
		syncCommitteeDuty := testingutils.TestingCommitteeDutyForSlot(phase0.Slot(pastHeight), nil, validatorsIndexList)
		attestationAndSyncCommitteeDuty := testingutils.TestingCommitteeDutyForSlot(phase0.Slot(pastHeight), validatorsIndexList, validatorsIndexList)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
			{
				Name: fmt.Sprintf("%v attestation (%s)", numValidators, version.String()),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap),
					attesterDuty,
					attestationMessages),
				Input: []interface{}{
					testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name: fmt.Sprintf("%v sync committee (%s)", numValidators, version.String()),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap),
					syncCommitteeDuty,
					syncCommitteeMessages),
				Input: []interface{}{
					testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name: fmt.Sprintf("%v attestation %v sync committee (%s)", numValidators, numValidators, version.String()),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap),
					attestationAndSyncCommitteeDuty,
					attestationAndSyncCommitteeMessages),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		}...)
	}

	return multiSpecTest
}
