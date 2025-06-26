package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapSyncCommittee tests mapping of BNRoleSyncCommittee.
func MapSyncCommittee() *DutySpecTest {
	return NewDutySpecTest(
		"map sync committee role",
		"Test mapping of BNRoleSyncCommittee",
		types.BNRoleSyncCommittee,
		types.RoleCommittee,
	)
}
