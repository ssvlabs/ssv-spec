package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorValidation tests a valid consensus data with AggregateAndProof
func Phase0AggregatorValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"phase0 aggregator valid",
		"Test validation of valid consensus data with Phase0 AggregateAndProof",
		*testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
		"",
	)
}

// ElectraAggregatorValidation tests a valid consensus data with AggregateAndProof
func ElectraAggregatorValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"electra aggregator valid",
		"Test validation of valid consensus data with Electra AggregateAndProof",
		*testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
		"",
	)
}
