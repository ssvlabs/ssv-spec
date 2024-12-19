package share

import (
	"github.com/ssvlabs/ssv-spec/types"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
	"github.com/stretchr/testify/require"
	reflect2 "reflect"
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
	decodedShare := &types.Share{}
	require.NoError(t, decodedShare.Decode(test.Data))

	// encode
	byts, err := decodedShare.Encode()
	require.NoError(t, err)
	require.EqualValues(t, test.Data, byts)

	// test root
	r, err := decodedShare.HashTreeRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedRoot, r)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}
