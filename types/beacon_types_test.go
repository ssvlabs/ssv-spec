package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapDutyToRunnerRole(t *testing.T) {
	testCases := []struct {
		name           string
		input          BeaconRole
		expectedOutput RunnerRole
	}{
		{
			name:           "Test Attester Role",
			input:          BNRoleAttester,
			expectedOutput: RoleCommittee,
		},
		{
			name:           "Test Proposer Role",
			input:          BNRoleProposer,
			expectedOutput: RoleProposer,
		},
		{
			name:           "Test Aggregator Role",
			input:          BNRoleAggregator,
			expectedOutput: RoleAggregator,
		},
		{
			name:           "Test Sync Committee Role",
			input:          BNRoleSyncCommittee,
			expectedOutput: RoleCommittee,
		},
		{
			name:           "Test Sync Committee Contribution Role",
			input:          BNRoleSyncCommitteeContribution,
			expectedOutput: RoleSyncCommitteeContribution,
		},
		{
			name:           "Test Validator Registration Role",
			input:          BNRoleValidatorRegistration,
			expectedOutput: RoleValidatorRegistration,
		},
		{
			name:           "Test Voluntary Exit Role",
			input:          BNRoleVoluntaryExit,
			expectedOutput: RoleVoluntaryExit,
		},
		{
			name:           "Test Unknown Role",
			input:          math.MaxInt32,
			expectedOutput: RoleUnknown,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := MapDutyToRunnerRole(tc.input)
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}
