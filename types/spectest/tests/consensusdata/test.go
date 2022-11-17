package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
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
