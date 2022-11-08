package qbft

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestController_Marshaling(t *testing.T) {
	c := testingControllerStruct

	byts, err := c.Encode()
	require.NoError(t, err)

	decoded := &Controller{}
	require.NoError(t, decoded.Decode(byts))

	bytsDecoded, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, bytsDecoded)
}
