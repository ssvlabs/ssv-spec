package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func Phase0AggregatorNoJustifications() *ValidatorConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return &ValidatorConsensusDataTest{
		Name:          "phase0 aggregator without justification",
		ConsensusData: *testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
	}
}

// ElectraAggregatorNoJustifications tests an invalid consensus data with no aggregator pre-consensus justifications
func ElectraAggregatorNoJustifications() *ValidatorConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return &ValidatorConsensusDataTest{
		Name:          "electra aggregator without justification",
		ConsensusData: *testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
	}
}
