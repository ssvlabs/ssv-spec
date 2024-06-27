package validatorconsensusdata

import (
	reflect2 "reflect"
	"testing"

	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type ValidatorConsensusDataTest struct {
	Name          string
	ConsensusData types.ValidatorConsensusData
	ExpectedError string
}

func (test *ValidatorConsensusDataTest) TestName() string {
	return "validatorconsensusdata " + test.Name
}

func (test *ValidatorConsensusDataTest) Run(t *testing.T) {

	err := test.ConsensusData.Validate()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}
