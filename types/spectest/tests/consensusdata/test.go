package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type EncodingSpecTest struct {
	Name string
	Data []byte
}

func (test *EncodingSpecTest) TestName() string {
	return test.Name
}

func (test *EncodingSpecTest) Run(t *testing.T) {
	a := types.ConsensusData{}
	require.NoError(t, a.Decode(test.Data))

	byts, err := a.Encode()
	require.NoError(t, err)
	require.EqualValues(t, test.Data, byts)
}
