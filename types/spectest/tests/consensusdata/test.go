package consensusdata

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type EncodingSpecTest struct {
	Name string
	Obj  *types.ConsensusData
}

func (test *EncodingSpecTest) TestName() string {
	return test.Name
}

func (test *EncodingSpecTest) Run(t *testing.T) {
	byts, err := test.Obj.Encode()
	require.NoError(t, err)

	a := types.ConsensusData{}
	require.NoError(t, a.Decode(byts))
	bytsDecoded, err := a.Encode()
	require.NoError(t, err)

	require.EqualValues(t, bytsDecoded, byts)
}

type ValidationSpecTest struct {
	Name        string
	Obj         *types.ConsensusData
	ExpectedErr string
}

func (test *ValidationSpecTest) TestName() string {
	return test.Name
}

func (test *ValidationSpecTest) Run(t *testing.T) {
	err := test.Obj.Validate()
	if len(test.ExpectedErr) > 0 {
		require.EqualError(t, err, test.ExpectedErr)
	} else {
		require.NoError(t, err)
	}
}
