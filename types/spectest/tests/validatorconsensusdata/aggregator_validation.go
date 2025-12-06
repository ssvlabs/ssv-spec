package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorValidation tests a valid consensus data with AggregateAndProof
func Phase0AggregatorValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"phase0 aggregator valid",
		testdoc.ValidatorConsensusDataTestPhase0AggregatorValidationDoc,
		*testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
		0,
	)
}

// ElectraAggregatorValidation tests a valid consensus data with AggregateAndProof
func ElectraAggregatorValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"electra aggregator valid",
		testdoc.ValidatorConsensusDataTestElectraAggregatorValidationDoc,
		*testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
		0,
	)
}
