package validatorconsensusdata

import "github.com/ssvlabs/ssv-spec/types/testingutils"

// InvalidAggregatorValidation tests an invalid consensus data with AggregateAndProof
func InvalidAggregatorValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData

	cd.DataSSZ = testingutils.TestingSyncCommitteeBlockRoot[:]

	return &ValidatorConsensusDataTest{
		Name:          "invalid aggregator data",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
