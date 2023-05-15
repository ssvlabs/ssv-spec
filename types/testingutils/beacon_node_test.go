package testingutils

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockRoot(t *testing.T) {
	r1, _ := TestingBeaconBlockV(spec.DataVersionBellatrix).Root()
	r2, _ := TestingBlindedBeaconBlockV(spec.DataVersionBellatrix).Root()
	require.EqualValues(t, r1, r2)
}
