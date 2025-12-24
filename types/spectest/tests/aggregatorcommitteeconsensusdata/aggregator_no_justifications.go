package aggregatorcommitteeconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func Phase0AggregatorNoJustifications() *AggregatorCommitteeConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return NewValidatorConsensusDataTest(
		"phase0 aggregator without justification",
		testdoc.AggregatorCommitteeConsensusDataTestPhase0AggregatorNoJustificationsDoc,
		*testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
		0,
	)
}

// ElectraAggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func ElectraAggregatorNoJustifications() *AggregatorCommitteeConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return NewValidatorConsensusDataTest(
		"electra aggregator without justification",
		testdoc.AggregatorCommitteeConsensusDataTestElectraAggregatorNoJustificationsDoc,
		*testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
		0,
	)
}
