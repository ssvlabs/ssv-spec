package consensusdata

import "github.com/bloxapp/ssv-spec/types/testingutils"

// InvalidSyncCommitteeContributionValidation tests an invalid consensus data with sync committee contrib.
func InvalidSyncCommitteeContributionValidation() *ConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cd.DataSSZ = testingutils.TestingAttestationDataBytes

	return &ConsensusDataTest{
		Name:          "invalid sync committee contribution",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: four",
	}
}
