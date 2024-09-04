package validation

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/types/testingutils"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
)

// PartialSigTypeMismatch tests the matching between the role and the signature type for partial signature messages
func PartialSigTypeMismatch() tests.SpecTest {

	expectedErr := validation.ErrPartialSignatureTypeRoleMismatch.Error()

	// Test cases: Role -> Partial signature type -> expectedError
	roleTypeCombinations := map[types.RunnerRole]map[types.PartialSigMsgType]string{
		types.RoleCommittee: {
			types.PostConsensusPartialSig:         "",
			types.RandaoPartialSig:                expectedErr,
			types.SelectionProofPartialSig:        expectedErr,
			types.ContributionProofs:              expectedErr,
			types.ValidatorRegistrationPartialSig: expectedErr,
			types.VoluntaryExitPartialSig:         expectedErr,
		},
		types.RoleProposer: {
			types.PostConsensusPartialSig:         "",
			types.RandaoPartialSig:                "",
			types.SelectionProofPartialSig:        expectedErr,
			types.ContributionProofs:              expectedErr,
			types.ValidatorRegistrationPartialSig: expectedErr,
			types.VoluntaryExitPartialSig:         expectedErr,
		},
		types.RoleAggregator: {
			types.PostConsensusPartialSig:         "",
			types.RandaoPartialSig:                expectedErr,
			types.SelectionProofPartialSig:        "",
			types.ContributionProofs:              expectedErr,
			types.ValidatorRegistrationPartialSig: expectedErr,
			types.VoluntaryExitPartialSig:         expectedErr,
		},
		types.RoleSyncCommitteeContribution: {
			types.PostConsensusPartialSig:         "",
			types.RandaoPartialSig:                expectedErr,
			types.SelectionProofPartialSig:        expectedErr,
			types.ContributionProofs:              "",
			types.ValidatorRegistrationPartialSig: expectedErr,
			types.VoluntaryExitPartialSig:         expectedErr,
		},
		types.RoleValidatorRegistration: {
			types.PostConsensusPartialSig:         expectedErr,
			types.RandaoPartialSig:                expectedErr,
			types.SelectionProofPartialSig:        expectedErr,
			types.ContributionProofs:              expectedErr,
			types.ValidatorRegistrationPartialSig: "",
			types.VoluntaryExitPartialSig:         expectedErr,
		},
		types.RoleVoluntaryExit: {
			types.PostConsensusPartialSig:         expectedErr,
			types.RandaoPartialSig:                expectedErr,
			types.SelectionProofPartialSig:        expectedErr,
			types.ContributionProofs:              expectedErr,
			types.ValidatorRegistrationPartialSig: expectedErr,
			types.VoluntaryExitPartialSig:         "",
		},
	}

	// Multi test creation
	multiTest := &MultiMessageValidationTest{
		Name:  "partial signature type mismatch",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for role, pSigTypeMap := range roleTypeCombinations {
		for pSigType, expectedError := range pSigTypeMap {
			multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
				Name:          fmt.Sprintf("role %v partial sig type %v", role, pSigType),
				Messages:      [][]byte{testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSignatureTypeAndRole(pSigType, role))},
				ExpectedError: expectedError,
			})
		}
	}

	return multiTest
}
