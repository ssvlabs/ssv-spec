package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorValidation tests a valid consensus data with AggregateAndProof
func Phase0AggregatorValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"phase0 aggregator valid",
		*testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
		"",
	)
}

// ElectraAggregatorValidation tests a valid consensus data with AggregateAndProof
func ElectraAggregatorValidation() *ValidatorConsensusDataTest {
	return NewValidatorConsensusDataTest(
		"electra aggregator valid",
		*testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
		"",
	)
}
