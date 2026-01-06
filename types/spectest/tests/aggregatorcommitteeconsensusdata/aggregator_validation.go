package aggregatorcommitteeconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorValidation tests a valid consensus data with AggregateAndProof
func Phase0AggregatorValidation() *AggregatorCommitteeConsensusDataTest {
	return NewAggregatorCommitteeConsensusDataTest(
		"phase0 aggregator valid",
		testdoc.AggregatorCommitteeConsensusDataTestPhase0AggregatorValidationDoc,
		*testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
		0,
	)
}

// ElectraAggregatorValidation tests a valid consensus data with AggregateAndProof
func ElectraAggregatorValidation() *AggregatorCommitteeConsensusDataTest {
	return NewAggregatorCommitteeConsensusDataTest(
		"electra aggregator valid",
		testdoc.AggregatorCommitteeConsensusDataTestElectraAggregatorValidationDoc,
		*testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
		0,
	)
}
