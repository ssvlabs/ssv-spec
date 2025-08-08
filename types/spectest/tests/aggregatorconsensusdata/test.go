package aggregatorconsensusdata

import (
	reflect2 "reflect"
	"testing"

	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/stretchr/testify/require"
)

type AggregatorConsensusDataTest struct {
	Name          string
	Type          string
	Documentation string
	ConsensusData types.AggregatorCommitteeConsensusData
	ExpectedError string
}

func (test *AggregatorConsensusDataTest) TestName() string {
	return "aggregatorconsensusdata " + test.Name
}

func (test *AggregatorConsensusDataTest) Run(t *testing.T) {

	err := test.ConsensusData.Validate()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewAggregatorConsensusDataTest(name, documentation string, consensusData types.AggregatorCommitteeConsensusData, expectedError string) *AggregatorConsensusDataTest {
	return &AggregatorConsensusDataTest{
		Name:          name,
		Type:          testdoc.AggregatorConsensusDataTestType,
		Documentation: documentation,
		ConsensusData: consensusData,
		ExpectedError: expectedError,
	}
}
