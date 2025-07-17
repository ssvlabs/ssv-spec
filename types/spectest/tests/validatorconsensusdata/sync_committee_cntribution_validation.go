package validatorconsensusdata

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// SyncCommitteeContributionValidation tests a valid consensus data with sync committee contrib.
func SyncCommitteeContributionValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"sync committee contribution valid",
		"Test validation of valid consensus data with sync committee contribution",
		*testingutils.TestSyncCommitteeContributionConsensusData,
		"",
	)
}
