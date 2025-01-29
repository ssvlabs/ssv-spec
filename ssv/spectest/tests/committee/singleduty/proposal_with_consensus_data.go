package committeesingleduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposalWithConsensusData sends a proposal message to a cluster runner with a ValidatorConsensusData (which is wrong since it accepts only BeaconVote objects)
func ProposalWithConsensusData() tests.SpecTest {

	ks := testingutils.TestingKeySetMap[phase0.ValidatorIndex(1)]
	msgID := testingutils.CommitteeMsgID(ks)
	numValidators := 30
	validatorsIndexList := testingutils.ValidatorIndexList(numValidators)
	ksMap := testingutils.KeySetMapForValidators(numValidators)
	shareMap := testingutils.ShareMapFromKeySetMap(ksMap)

	expectedError := "failed processing consensus message: could not process msg: invalid signed message: proposal not justified: proposal fullData invalid: failed decoding beacon vote: incorrect size"

	multiSpecTest := &committee.MultiCommitteeSpecTest{
		Name:  "proposal with consensus data",
		Tests: []*committee.CommitteeSpecTest{},
	}

	for _, version := range testingutils.SupportedAttestationVersions {

		slot := testingutils.TestingDutySlotV(version)
		height := qbft.Height(slot)

		multiSpecTest.Tests = append(multiSpecTest.Tests, []*committee.CommitteeSpecTest{
			{
				Name:      fmt.Sprintf("%v attestation (%s)", numValidators, version.String()),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingAttesterDutyForValidators(version, validatorsIndexList),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestAttesterConsensusDataByts,
						height),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("%v sync committee (%s)", numValidators, version.String()),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingSyncCommitteeDutyForValidators(version, validatorsIndexList),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestSyncCommitteeConsensusDataByts,
						height),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
			{
				Name:      fmt.Sprintf("%v attestations %v sync committees (%s)", numValidators, numValidators, version.String()),
				Committee: testingutils.BaseCommitteeWithCreatorFieldsFromRunner(ksMap, testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)),
				Input: []interface{}{
					testingutils.TestingCommitteeDuty(validatorsIndexList, validatorsIndexList, version),
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), msgID, testingutils.TestAttesterConsensusDataByts,
						height),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedError,
			},
		}...)
	}

	return multiSpecTest
}
