package validatorconsensusdata

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// InvalidSyncCommitteeContributionValidation tests an invalid consensus data with sync committee contrib.
func InvalidSyncCommitteeContributionValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cd.DataSSZ = testingutils.TestingAttestationDataBytes

	return &ValidatorConsensusDataTest{
		Name:          "invalid sync committee contribution",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: four",
	}
}
