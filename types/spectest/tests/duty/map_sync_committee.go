package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapSyncCommittee tests mapping of BNRoleSyncCommittee.
func MapSyncCommittee() *DutySpecTest {
	return NewDutySpecTest(
		"map sync committee role",
		testdoc.MapSyncCommitteeTestDoc,
		types.BNRoleSyncCommittee,
		types.RoleCommittee,
	)
}
