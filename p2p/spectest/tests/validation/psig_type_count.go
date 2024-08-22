package validation

import (
	"fmt"

	"github.com/ssvlabs/ssv-spec/types/testingutils"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
)

// PartialSigTypeCount tests sending two partial signature messages and trigerring the message count rule
func PartialSigTypeCount() tests.SpecTest {

	// Function to generate to 2 partial signature messages with different signing roots
	messages := func(role types.RunnerRole, pSigType types.PartialSigMsgType) [][]byte {
		return [][]byte{testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSignatureTypeRoleAndRoot(pSigType, role, [32]byte{1})),
			testingutils.EncodeMessage(testingutils.PartialSignatureMsgForSignatureTypeRoleAndRoot(pSigType, role, [32]byte{2})),
		}
	}

	// Test cases: possible combinations between roles and signature types
	testCases := map[types.RunnerRole]map[types.PartialSigMsgType]struct{}{
		types.RoleCommittee: {
			types.PostConsensusPartialSig: {},
		},
		types.RoleProposer: {
			types.RandaoPartialSig:        {},
			types.PostConsensusPartialSig: {},
		},
		types.RoleAggregator: {
			types.SelectionProofPartialSig: {},
			types.PostConsensusPartialSig:  {},
		},
		types.RoleSyncCommitteeContribution: {
			types.ContributionProofs:      {},
			types.PostConsensusPartialSig: {},
		},
		types.RoleValidatorRegistration: {
			types.ValidatorRegistrationPartialSig: {},
		},
		types.RoleVoluntaryExit: {
			types.VoluntaryExitPartialSig: {},
		},
	}

	expectedErr := validation.ErrInvalidPartialSignatureTypeCount.Error()

	// Multi test creation
	multiTests := &MultiMessageValidationTest{
		Name:  "partial signature type count",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for role, pSigTypeMap := range testCases {
		for pSigType := range pSigTypeMap {
			multiTests.Tests = append(multiTests.Tests, &MessageValidationTest{
				Name:          fmt.Sprintf("role %v partial signature type %v", role, pSigType),
				Messages:      messages(role, pSigType),
				ExpectedError: expectedErr,
			})
		}
	}

	return multiTests
}
