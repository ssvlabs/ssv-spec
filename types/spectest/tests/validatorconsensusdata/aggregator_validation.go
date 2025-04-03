package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Phase0AggregatorValidation tests a valid consensus data with AggregateAndProof
func Phase0AggregatorValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "phase0 aggregator valid",
		ConsensusData: *testingutils.TestAggregatorConsensusData(spec.DataVersionPhase0),
	}
}

// ElectraAggregatorValidation tests a valid consensus data with AggregateAndProof
func ElectraAggregatorValidation() *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:          "electra aggregator valid",
		ConsensusData: *testingutils.TestAggregatorConsensusData(spec.DataVersionElectra),
	}
}
