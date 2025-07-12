package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapProposer tests mapping of BNRoleProposer.
func MapProposer() *DutySpecTest {
	return NewDutySpecTest(
		"map proposer role",
		"Test mapping of BNRoleProposer",
		types.BNRoleProposer,
		types.RoleProposer,
	)
}
