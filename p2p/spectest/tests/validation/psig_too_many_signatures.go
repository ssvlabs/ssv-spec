package validation

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigTooManySignatures tests partial singature messages with many signatures
func PartialSigTooManySignatures() tests.SpecTest {

	// Number of validators of the default CommitteeID
	numCommitteeValidators := len(testingutils.NewTestingNetworkDataFetcher().Committees[testingutils.TestingCommitteeID].Validators)

	expectedErr := validation.ErrTooManyPartialSignatureMessages.Error()

	// Test cases: Role -> Number of signatures
	testCases := map[types.RunnerRole]map[int]string{
		types.RoleCommittee: {
			1:                            "",
			numCommitteeValidators:       "",
			numCommitteeValidators + 1:   "",
			2 * numCommitteeValidators:   "",
			2*numCommitteeValidators + 1: expectedErr,
		},
		types.RoleProposer: {
			1: "",
			2: expectedErr,
		},
		types.RoleAggregator: {
			1: "",
			2: expectedErr,
		},
		types.RoleSyncCommitteeContribution: {
			1:                            "",
			validation.MaxSignatures:     "",
			validation.MaxSignatures + 1: expectedErr,
		},
		types.RoleValidatorRegistration: {
			1: "",
			2: expectedErr,
		},
		types.RoleVoluntaryExit: {
			1: "",
			2: expectedErr,
		},
	}

	// Multi test creation
	multiTests := &MultiMessageValidationTest{
		Name:  "partial signature too many signatures",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for role, numSignaturesMap := range testCases {
		for numSignatures, testCaseErr := range numSignaturesMap {
			multiTests.Tests = append(multiTests.Tests, &MessageValidationTest{
				Name:          fmt.Sprintf("role %v with %v signatures", role, numSignatures),
				Messages:      [][]byte{testingutils.EncodeMessage(testingutils.PartialSignatureMsgForNumSignatures(numSignatures, role, numCommitteeValidators))},
				ExpectedError: testCaseErr,
			})
		}
	}

	return multiTests
}
