package testingutils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBeaconBlockRoot(t *testing.T) {
	r1, _ := TestingBeaconBlock.HashTreeRoot()
	r2, _ := TestingBlindedBeaconBlock.HashTreeRoot()
	require.EqualValues(t, r1, r2)
}
