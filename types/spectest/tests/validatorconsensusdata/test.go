package validatorconsensusdata

import (
	reflect2 "reflect"
	"testing"

	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/errcodes"
	"github.com/stretchr/testify/require"
)

type ValidatorConsensusDataTest struct {
	Name              string
	Type              string
	Documentation     string
	ConsensusData     types.ValidatorConsensusData
	ExpectedErrorCode errcodes.Code
}

func (test *ValidatorConsensusDataTest) TestName() string {
	return "validatorconsensusdata " + test.Name
}

func (test *ValidatorConsensusDataTest) Run(t *testing.T) {

	err := test.ConsensusData.Validate()

	if test.ExpectedErrorCode != 0 {
		require.Equal(t, test.ExpectedErrorCode, errcodes.FromError(err))
	} else {
		require.NoError(t, err)
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewValidatorConsensusDataTest(name, documentation string, consensusData types.ValidatorConsensusData, expectedErrorCode errcodes.Code) *ValidatorConsensusDataTest {
	return &ValidatorConsensusDataTest{
		Name:              name,
		Type:              testdoc.ValidatorConsensusDataTestType,
		Documentation:     documentation,
		ConsensusData:     consensusData,
		ExpectedErrorCode: expectedErrorCode,
	}
}
