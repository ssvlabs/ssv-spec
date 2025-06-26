package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapSyncCommitteeContribution tests mapping of BNRoleSyncCommitteeContribution.
func MapSyncCommitteeContribution() *DutySpecTest {
	return NewDutySpecTest(
		"map sync committee contribution role",
		"Test mapping of BNRoleSyncCommitteeContribution",
		types.BNRoleSyncCommitteeContribution,
		types.RoleSyncCommitteeContribution,
	)
}
