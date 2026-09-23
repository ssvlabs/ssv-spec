package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapProposerPreferences tests mapping of BNRoleProposerPreferences.
func MapProposerPreferences() *DutySpecTest {
	return NewDutySpecTest(
		"map proposer preferences role",
		testdoc.MapProposerPreferencesTestDoc,
		types.BNRoleProposerPreferences,
		types.RoleProposerPreferences,
	)
}
