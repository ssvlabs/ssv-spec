package signedssvmsg

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type EncodingTest struct {
	Name string
	Type string
	Data []byte
}

func (test *EncodingTest) TestName() string {
	return test.Name
}

func (test *EncodingTest) Run(t *testing.T) {
	// decode
	decoded := &types.SignedSSVMessage{}
	require.NoError(t, decoded.Decode(test.Data))

	// encode
	byts, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, test.Data, byts)
}

func NewEncodingTest(name string, data []byte) *EncodingTest {
	return &EncodingTest{
		Name: name,
		Type: "Signed SSV message encoding",
		Data: data,
	}
}
