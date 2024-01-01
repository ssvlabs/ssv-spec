package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AggregatorNoJustifications tests consensus data with no aggregator pre-consensus justifications
func AggregatorNoJustifications() tests.SpecTest {
	cd := testingutils.TestAggregatorConsensusData

	byts, err := cd.Encode()
	if err != nil {
		panic(err.Error())
	}
	root, err := cd.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}

	return &EncodingTest{
		Name:         "consensusdata aggregator no justifications encoding",
		Data:         byts,
		ExpectedRoot: root,
	}
}
