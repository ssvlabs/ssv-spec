package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
)

// AssertErrorCode asserts the error against expected error code.
func AssertErrorCode(t *testing.T, wantErrCode int, err error) {
	if wantErrCode == 0 {
		require.NoError(t, err)
		return
	}

	e := &types.Error{}
	if !errors.As(err, &e) {
		require.Failf(t, "unknown error", "%+v", err)
	}
	require.Equal(t, wantErrCode, e.Code)
}
