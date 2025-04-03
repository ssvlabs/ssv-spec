package committeesingleduty

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PastMessageDutyDoesNotExist tests a valid proposal past msg for a duty that doesn't exist
func PastMessageDutyDoesNotExist() tests.SpecTest {

	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	ks := testingutils.Testing4SharesSet()

	decidedValue := testingutils.TestBeaconVoteByts
	msgID := testingutils.CommitteeMsgID(ks)
	pastHeight := qbft.Height(10)

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

	expectedError := "no runner found for message's slot"

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "past msg duty does not exist",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {
		multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
			{
				Name:      fmt.Sprintf("%v attestation (%s)", numValidators, version.String()),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("%v sync committee (%s)", numValidators, version.String()),
				Committee: testingutils.BaseCommittee(ksMap),
				Input: []interface{}{
					testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					pastProposalMsgF(),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("%v attestation %v sync committee (%s)", numValidators, numValidators, version.String()),
				Committee: testingutils.BaseCommittee(ksMap),
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
