package aggregatorcommitteeconsensusdata

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

type AggregatorCommitteeConsensusDataTest struct {
	Name                 string
	Type                 string
	Documentation        string
	AggCommConsensusData types.AggregatorCommitteeConsensusData
	ExpectedErrorCode    int
}

func (test *AggregatorCommitteeConsensusDataTest) TestName() string {
	return "aggregatorcommitteeconsensusdata " + test.Name
}

func (test *AggregatorCommitteeConsensusDataTest) Run(t *testing.T) {

	err := test.AggCommConsensusData.Validate()
	tests.AssertErrorCode(t, test.ExpectedErrorCode, err)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewValidatorConsensusDataTest(name, documentation string, cd types.AggregatorCommitteeConsensusData, expectedErrorCode int) *AggregatorCommitteeConsensusDataTest {
	return &AggregatorCommitteeConsensusDataTest{
		Name:                 name,
		Type:                 testdoc.AggregatorCommitteeConsensusDataTestType,
		Documentation:        documentation,
		AggCommConsensusData: cd,
		ExpectedErrorCode:    expectedErrorCode,
	}
}
