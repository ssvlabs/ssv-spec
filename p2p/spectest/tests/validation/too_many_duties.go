package validation

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TooManyDutiesWithConsensusMessages tests consensus messages with different heights to indicate a sequence of duties
func TooManyDutiesWithConsensusMessages() tests.SpecTest {

	// Functino to create sequence of consensus messages given a number of duties
	consensusMsgsForRole := func(role types.RunnerRole, numDuties int, committeeID types.CommitteeID) [][]byte {

		msgID := testingutils.MessageIDForRoleAndCommitteeID(role, committeeID)

		ret := [][]byte{}
		for slot := 0; slot < numDuties; slot++ {
			ret = append(ret, testingutils.EncodeMessage(testingutils.ConsensusMsgForSlot(phase0.Slot(slot), msgID, testingutils.KeySetForCommitteeID[committeeID])))
		}
		return ret
	}

	expectedErr := validation.ErrTooManyDutiesPerEpoch.Error()

	// Test cases: Role -> CommitteeID -> Number of duties per epoch -> ExpectedError
	// the CommitteeID is used here since, for the RoleCommittee, it influences the number of duties per epoch
	roleTestCases := map[types.RunnerRole]map[types.CommitteeID]map[int]string{
		types.RoleAggregator: {
			testingutils.TestingCommitteeID: {
				2: "",
				3: expectedErr,
			},
		},
		types.RoleCommittee: {
			testingutils.TestingCommitteeID: { // No validators in sync committee, as defined in the testing MessageValidator.NetworkDataFetcher. Has 10 validators
				20: "",
				21: expectedErr,
			},
			testingutils.TestingCommitteeIDWithSyncCommitteeDuty: { // Has a validator in sync committee, as defined in the testing MessageValidator.NetworkDataFetcher
				1:  "",
				32: "",
			},
		},
	}

	// Create multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "too many duties with consensus messages",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for role, committeeIDMap := range roleTestCases {
		for committeeID, numDutiesMap := range committeeIDMap {
			for numDuties, testCaseErr := range numDutiesMap {
				multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
					Name:          fmt.Sprintf("role %v with %v duties", role, numDuties),
					Messages:      consensusMsgsForRole(role, numDuties, committeeID),
					ExpectedError: testCaseErr,
				})
			}
		}
	}

	return multiTest
}

// TooManyDutiesWithPartialSignatureMessages tests partial signature messages with different heights to indicate a sequence of duties
func TooManyDutiesWithPartialSignatureMessages() tests.SpecTest {

	// Create sequence of messagess given a number of duties for non committee roles
	msgsForRole := func(role types.RunnerRole, numDuties int, committeeID types.CommitteeID) [][]byte {

		msgID := testingutils.MessageIDForRoleAndCommitteeID(role, committeeID)

		ret := [][]byte{}
		for slot := 0; slot < numDuties; slot++ {
			ret = append(ret, testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSlot(phase0.Slot(slot), msgID, testingutils.ValidValidatorIndexForCommitteeID[committeeID], testingutils.KeySetForCommitteeID[committeeID])))
		}
		return ret
	}

	expectedErr := validation.ErrTooManyDutiesPerEpoch.Error()

	// Test cases: Role -> CommitteeID -> Number of duties per epoch -> ExpectedError
	// the CommitteeID is used here since, for the RoleCommittee, it influences the number of duties per epoch
	roleTestCases := map[types.RunnerRole]map[types.CommitteeID]map[int]string{
		types.RoleAggregator: {
			testingutils.TestingCommitteeID: {
				2: "",
				3: expectedErr,
			},
		},
		types.RoleValidatorRegistration: {
			testingutils.TestingCommitteeID: {
				2: "",
				3: expectedErr,
			},
		},
		types.RoleVoluntaryExit: {
			testingutils.TestingCommitteeID: {
				2: "",
				3: expectedErr,
			},
		},
		types.RoleCommittee: {
			testingutils.TestingCommitteeID: { // No validators in sync committee, as defined in the testing MessageValidator.NetworkDataFetcher. Has 10 validators
				20: "",
				21: expectedErr,
			},
			testingutils.TestingCommitteeIDWithSyncCommitteeDuty: { // Has a validator in sync committee, as defined in the testing MessageValidator.NetworkDataFetcher
				1:  "",
				32: "",
			},
		},
	}

	// Construct multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "too many duties with partial signature messages",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for role, committeeIDMap := range roleTestCases {
		for committeeID, numDutiesMap := range committeeIDMap {
			for numDuties, testCaseErr := range numDutiesMap {
				multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
					Name:          fmt.Sprintf("role %v with %v duties", role, numDuties),
					Messages:      msgsForRole(role, numDuties, committeeID),
					ExpectedError: testCaseErr,
				})
			}
		}
	}

	return multiTest
}
