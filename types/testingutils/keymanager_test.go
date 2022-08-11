package testingutils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTestingTimer(t *testing.T) {
	r, _ := TestingAttestationData.HashTreeRoot()
	require.EqualValues(t, r[:], TestingAttestationRoot)
}
