package payloadattestationdata

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
	"github.com/stretchr/testify/require"
)

type EncodingTest struct {
	Name          string
	Type          string
	Documentation string
	Data          []byte
	ExpectedRoot  [32]byte
	// DecodeFails marks data that must not decode as a PayloadAttestationData. ExpectedRoot is unused then.
	DecodeFails bool `json:",omitempty"`
}

func (test *EncodingTest) TestName() string {
	return test.Name
}

func (test *EncodingTest) Run(t *testing.T) {
	// decode
	decoded := &gloas.PayloadAttestationData{}
	if test.DecodeFails {
		require.Error(t, decoded.UnmarshalSSZ(test.Data))
		comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
		return
	}
	require.NoError(t, decoded.UnmarshalSSZ(test.Data))

	// encode
	byts, err := decoded.MarshalSSZ()
	require.NoError(t, err)
	require.EqualValues(t, test.Data, byts)

	// test root
	r, err := decoded.HashTreeRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ExpectedRoot, r)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewEncodingTest(name, documentation string, data []byte, expectedRoot [32]byte) *EncodingTest {
	return &EncodingTest{
		Name:          name,
		Type:          testdoc.PayloadAttestationDataEncodingSpecTestType,
		Documentation: documentation,
		Data:          data,
		ExpectedRoot:  expectedRoot,
	}
}
