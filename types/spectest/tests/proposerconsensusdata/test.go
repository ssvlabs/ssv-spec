package proposerconsensusdata

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

type ProposerConsensusDataTest struct {
	Name              string
	Type              string
	Documentation     string
	ConsensusData     types.ProposerConsensusData
	ExpectedErrorCode int
}

func (test *ProposerConsensusDataTest) TestName() string {
	return "proposerconsensusdata " + test.Name
}

func (test *ProposerConsensusDataTest) Run(t *testing.T) {

	err := test.ConsensusData.Validate()
	tests.AssertErrorCode(t, test.ExpectedErrorCode, err)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewProposerConsensusDataTest(name, documentation string, consensusData types.ProposerConsensusData, expectedErrorCode int) *ProposerConsensusDataTest {
	return &ProposerConsensusDataTest{
		Name:              name,
		Type:              testdoc.ProposerConsensusDataTestType,
		Documentation:     documentation,
		ConsensusData:     consensusData,
		ExpectedErrorCode: expectedErrorCode,
	}
}
