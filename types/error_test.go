package types

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	someErr := fmt.Errorf("some error")
	require.False(t, errors.Is(someErr, &Error{}))

	newErr := NewError(UnmarshalSSZErrorCode, someErr.Error())
	require.True(t, errors.Is(newErr, &Error{}))
	wrappedErr := WrapError(UnmarshalSSZErrorCode, someErr)
	require.True(t, errors.Is(wrappedErr, &Error{}))

	rErr := &Error{}
	require.False(t, errors.As(someErr, &rErr))
	require.True(t, errors.As(newErr, &rErr))
	require.True(t, errors.As(wrappedErr, &rErr))
}
