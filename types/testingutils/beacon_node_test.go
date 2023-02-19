package testingutils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBeaconBlockRoot(t *testing.T) {
	r1, _ := TestingBellatrixBeaconBlock.HashTreeRoot()
	r2, _ := TestingBellatrixBlindedBeaconBlock.HashTreeRoot()
	require.EqualValues(t, r1, r2)
}
