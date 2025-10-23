package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0InvalidAggregatorValidation tests an invalid consensus data with AggregateAndProof
func Phase0InvalidAggregatorValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0)

	cd.DataSSZ = testingutils.TestingSyncCommitteeBlockRoot[:]

	return NewValidatorConsensusDataTest(
		"invalid phase0 aggregator data",
		testdoc.ValidatorConsensusDataTestInvalidPhase0AggregatorDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}

// ElectraInvalidAggregatorValidation tests an invalid consensus data with AggregateAndProof
func ElectraInvalidAggregatorValidation() *ValidatorConsensusDataTest {

	cd := testingutils.TestAggregatorConsensusData(spec.DataVersionElectra)

	cd.DataSSZ = testingutils.TestingSyncCommitteeBlockRoot[:]

	return NewValidatorConsensusDataTest(
		"invalid electra aggregator data",
		testdoc.ValidatorConsensusDataTestInvalidElectraAggregatorDoc,
		*cd,
		types.UnmarshalSSZErrorCode,
	)
}
