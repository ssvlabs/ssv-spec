package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidSyncCommitteeContributionValidation tests an invalid consensus data with sync committee contrib.
func InvalidSyncCommitteeContributionValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cd.DataSSZ = testingutils.TestingAttestationDataBytes(spec.DataVersionPhase0)

	return NewValidatorConsensusDataTest(
		"invalid sync committee contribution",
		"Test validation error for invalid consensus data with sync committee contribution using incorrect data",
		*cd,
		"could not unmarshal ssz: incorrect end of offset: 12 0",
	)
}
