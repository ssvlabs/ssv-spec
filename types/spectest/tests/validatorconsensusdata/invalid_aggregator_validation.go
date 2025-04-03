package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0InvalidAggregatorValidation tests an invalid consensus data with AggregateAndProof
func Phase0InvalidAggregatorValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)

	cd.DataSSZ = testingutils.TestingSyncCommitteeBlockRoot[:]

	return &ValidatorConsensusDataTest{
		Name:          "invalid phase0 aggregator data",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}

// ElectraInvalidAggregatorValidation tests an invalid consensus data with AggregateAndProof
func ElectraInvalidAggregatorValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionElectra)

	cd.DataSSZ = testingutils.TestingSyncCommitteeBlockRoot[:]

	return &ValidatorConsensusDataTest{
		Name:          "invalid electra aggregator data",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
