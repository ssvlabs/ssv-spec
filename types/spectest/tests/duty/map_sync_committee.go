package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapSyncCommittee tests mapping of BNRoleSyncCommittee.
func MapSyncCommittee() *DutySpecTest {
	return NewDutySpecTest(
		"map sync committee role",
		types.BNRoleSyncCommittee,
		types.RoleCommittee,
	)
}
