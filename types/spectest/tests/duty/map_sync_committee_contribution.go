package duty

import "github.com/ssvlabs/ssv-spec/types"

// MapSyncCommitteeContribution tests mapping of BNRoleSyncCommitteeContribution.
func MapSyncCommitteeContribution() *DutySpecTest {
	return &DutySpecTest{
		Name:       "map sync committee contribution role",
		BeaconRole: types.BNRoleSyncCommitteeContribution,
		RunnerRole: types.RoleSyncCommitteeContribution,
	}
}
