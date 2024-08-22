package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundTooHigh tests a consensus message with round above the allowed limit
func RoundTooHigh() tests.SpecTest {

	expectedErr := validation.ErrRoundTooHigh.Error()

	// Test cases: Role -> Round -> Expected error
	testCases := map[types.RunnerRole]map[qbft.Round]string{
		types.RoleProposer: {
			6: "",
			7: expectedErr,
		},
		types.RoleSyncCommitteeContribution: {
			6: "",
			7: expectedErr,
		},
		types.RoleCommittee: {
			12: "",
			13: expectedErr,
		},
		types.RoleAggregator: {
			12: "",
			13: expectedErr,
		},
	}

	// Contruct multi test
	multiTest := &MultiMessageValidationTest{
		Name:  "round too high",
		Tests: make([]*MessageValidationTest, 0),
	}

	// Add test cases
	for role, roundMap := range testCases {
		for round, testCaseErr := range roundMap {
			multiTest.Tests = append(multiTest.Tests, &MessageValidationTest{
				Name:          "round too high for proposer",
				Messages:      [][]byte{testingutils.EncodeMessage(testingutils.ConsensusMessageForRound(round, testingutils.MessageIDForRole(role)))},
				ExpectedError: testCaseErr,
				ReceivedAt:    testingutils.ReceivedAtForRound(round - 1),
			})
		}
	}

	return multiTest
}
