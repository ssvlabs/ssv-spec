package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapSyncCommittee tests mapping of BNRoleSyncCommittee.
func MapSyncCommittee() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map sync committee role",
		BeaconRole: types.BNRoleSyncCommittee,
		RunnerRole: types.RoleCommittee,
	}
}
