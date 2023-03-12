package partialsigmessage

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type EncodingTest struct {
	Name         string
	Data         []byte
	ExpectedRoot [32]byte
}

func (test *EncodingTest) TestName() string {
	return test.Name
}

func (test *EncodingTest) Run(t *testing.T) {
	// decode
	decoded := &types.SignedPartialSignatureMessage{}
	require.NoError(t, decoded.Decode(test.Data))

	// encode
	byts, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, test.Data, byts)

	// test root
	r, err := decoded.GetRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedRoot, r)
}
