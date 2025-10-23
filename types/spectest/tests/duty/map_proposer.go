package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapProposer tests mapping of BNRoleProposer.
func MapProposer() *DutySpecTest {
	return NewDutySpecTest(
		"map proposer role",
		testdoc.MapProposerTestDoc,
		types.BNRoleProposer,
		types.RoleProposer,
	)
}
