package validatorconsensusdata

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// SyncCommitteeContributionValidation tests a valid consensus data with sync committee contrib.
func SyncCommitteeContributionValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "sync committee contribution valid",
		ConsensusData: *testingutils.TestSyncCommitteeContributionConsensusData,
	}

}
