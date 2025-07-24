package signedssvmsg

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/stretchr/testify/require"
)

type EncodingTest struct {
	Name          string
	Type          string
	Documentation string
	Data          []byte
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

func NewEncodingTest(name string, documentation string, data []byte) *EncodingTest {
	return &EncodingTest{
		Name:          name,
		Type:          testdoc.SignedSSVMessageEncodingTestType,
		Documentation: documentation,
		Data:          data,
	}
}
