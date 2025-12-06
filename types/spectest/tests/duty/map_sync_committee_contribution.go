package duty

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// MapSyncCommitteeContribution tests mapping of BNRoleSyncCommitteeContribution.
func MapSyncCommitteeContribution() *DutySpecTest {
	return NewDutySpecTest(
		"map sync committee contribution role",
		testdoc.MapSyncCommitteeContributionTestDoc,
		types.BNRoleSyncCommitteeContribution,
		types.RoleSyncCommitteeContribution,
	)
}
