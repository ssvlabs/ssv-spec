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

// PastMessageDutyNotFinished tests a valid proposal past msg for a duty that didnt finish
func PastMessageDutyNotFinished() tests.SpecTest {

	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	ks := testingutils.Testing4SharesSet()

	decidedValue := testingutils.TestBeaconVoteByts
	msgID := testingutils.CommitteeMsgID(ks)

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "past msg duty not finished",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		pastHeight := qbft.Height(testingutils.TestingDutySlotV(version) - 2)

		bumpHeight := func(c *ssv.Committee, previousDuty types.Duty) *ssv.Committee {

			err := c.StartDuty(previousDuty.(*types.CommitteeDuty))
			if err != nil {
				panic(err)
			}

			decidingMsgs := []*types.SignedSSVMessage{
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
			for _, msg := range decidingMsgs {
				err := c.ProcessMessage(msg)
				if err != nil {
					panic(err)
				}
			}

			// Erase broadcasted messages due to test setup
			c.Runners[previousDuty.DutySlot()].GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs = make([]*types.SignedSSVMessage, 0)

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

		expectedError := "failed processing consensus message: not processing consensus message since consensus has already finished"

		attesterDuty := testingutils.TestingCommitteeDutyForSlot(phase0.Slot(pastHeight), validatorsIndexList, nil)
		syncCommitteeDuty := testingutils.TestingCommitteeDutyForSlot(phase0.Slot(pastHeight), nil, validatorsIndexList)
		attestationAndSyncCommitteeDuty := testingutils.TestingCommitteeDutyForSlot(phase0.Slot(pastHeight), validatorsIndexList, validatorsIndexList)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
			{
				Name:      fmt.Sprintf("%v attestation (%s)", numValidators, version.String()),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap), attesterDuty),
				Input: []interface{}{
					testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("%v sync committee (%s)", numValidators, version.String()),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap), syncCommitteeDuty),
				Input: []interface{}{
					testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("%v attestation %v sync committee (%s)", numValidators, numValidators, version.String()),
				Committee: bumpHeight(testingutils.BaseCommittee(ksMap), attestationAndSyncCommitteeDuty),
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
